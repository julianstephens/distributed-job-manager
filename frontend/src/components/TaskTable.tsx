import { Tooltip } from "@/components/ui/tooltip";
import type { Task } from "@/lib/api/aliases";
import { $api } from "@/lib/api/client";
import { convertUnixToDate, TABLE_PAGE_SIZE } from "@/lib/utils";
import { TaskStatus } from "@/types";
import {
  Box,
  ButtonGroup,
  createListCollection,
  Flex,
  Icon,
  IconButton,
  Pagination,
  Portal,
  SegmentGroup,
  Select,
  Table,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { FaX } from "react-icons/fa6";
import { LuChevronLeft, LuChevronRight } from "react-icons/lu";

export const TaskTable = () => {
  const { data, error, isLoading } = $api.useQuery("get", "/tasks");

  const sortFilters = createListCollection({
    items: [
      { label: "Title", value: "title" },
      { label: "Created At", value: "createdAt" },
      { label: "Update At", value: "updatedAt" },
      { label: "Status", value: "status" },
    ],
  });
  const sortDirections = createListCollection({
    items: [
      { label: "Ascending", value: "asc" },
      { label: "Descending", value: "desc" },
    ],
  });
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");
  const [statusFilter, setStatusFilter] = useState<string | null>(null);
  const [filteredData, setFilteredData] = useState<Task[] | null>(null);
  const [page, setPage] = useState(1);

  useEffect(() => {
    if (!isLoading && !error && data && data.data) {
      setFilteredData(data.data);
    }
  }, [data, isLoading, error]);

  const updateStatusFilter = (status: string | null) => {
    setStatusFilter(status);
  };

  const updateSort = (sortBy: string | null) => {
    if (!filteredData) return;
    if (!sortBy) {
      setFilteredData(data?.data || null);
      return;
    }
    const sortedData = [...filteredData].sort((a, b) => {
      if (sortBy === "title") {
        return a.title.localeCompare(b.title);
      } else if (sortBy === "createdAt") {
        return (
          new Date(a.createdAt! * 1000).getTime() -
          new Date(b.createdAt! * 1000).getTime()
        );
      } else if (sortBy === "updatedAt") {
        return (
          new Date(a.updatedAt! * 1000).getTime() -
          new Date(b.updatedAt! * 1000).getTime()
        );
      } else if (sortBy === "status") {
        return a.status - b.status;
      }
      return 0;
    });
    setFilteredData(sortedData);
  };

  const updateSortDirection = (direction: "asc" | "desc") => {
    if (!filteredData) return;
    const sortedData = [...filteredData].reverse();
    setFilteredData(sortedData);
    setSortDirection(direction);
  };

  const getPagedData = () => {
    if (!filteredData) return null;
    const startIndex = (page - 1) * TABLE_PAGE_SIZE;
    const endIndex = startIndex + TABLE_PAGE_SIZE;
    return filteredData.slice(startIndex, endIndex);
  };

  return (
    <Flex direction="column">
      <Flex id="filterBar" align="center" justify="space-between" mb="5">
        <Flex w="1/2" align="center" gap="2">
          <Select.Root
            collection={sortFilters}
            maxW="2/5"
            multiple={false}
            onValueChange={({ value }) => {
              updateSort(value[0] || null);
            }}
          >
            <Select.HiddenSelect />
            <Select.Control>
              <Select.Trigger>
                <Select.ValueText placeholder="Sort by" />
              </Select.Trigger>
              <Select.IndicatorGroup>
                <Select.ClearTrigger />
                <Select.Indicator />
              </Select.IndicatorGroup>
            </Select.Control>
            <Portal>
              <Select.Positioner>
                <Select.Content>
                  {sortFilters.items.map((item) => (
                    <Select.Item key={item.value} item={item}>
                      <Select.ItemText>{item.label}</Select.ItemText>
                      <Select.ItemIndicator />
                    </Select.Item>
                  ))}
                </Select.Content>
              </Select.Positioner>
            </Portal>
          </Select.Root>
          <SegmentGroup.Root
            w="fit"
            my="5"
            value={sortDirection}
            onValueChange={({ value }) => {
              updateSortDirection(value as "asc" | "desc");
            }}
            disabled={!filteredData || filteredData.length === 0}
          >
            <SegmentGroup.Items
              items={sortDirections.items.map((item) => item.value)}
              disabled={isLoading}
              cursor="pointer"
            />
            <SegmentGroup.Indicator />
          </SegmentGroup.Root>
        </Flex>
        <Flex w="1/2" align="center" justify="end" gap="2">
          <SegmentGroup.Root
            w="fit"
            my="5"
            value={statusFilter}
            onValueChange={({ value }) => {
              updateStatusFilter(value === statusFilter ? null : value);
            }}
          >
            <SegmentGroup.Items
              items={Object.values(TaskStatus)}
              disabled={isLoading}
              cursor="pointer"
            />
            <SegmentGroup.Indicator />
          </SegmentGroup.Root>
          <Tooltip content="Clear Filter" disabled={!statusFilter}>
            <IconButton
              size="xs"
              variant="ghost"
              rounded="full"
              colorPalette="red"
              onClick={() => updateStatusFilter(null)}
              disabled={!statusFilter}
            >
              <Icon size="xs">
                <FaX />
              </Icon>
            </IconButton>
          </Tooltip>
        </Flex>
      </Flex>
      <Box width="7/12" mx="auto">
        <Table.Root interactive>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeader>Task ID</Table.ColumnHeader>
              <Table.ColumnHeader>Title</Table.ColumnHeader>
              <Table.ColumnHeader>Status</Table.ColumnHeader>
              <Table.ColumnHeader>Created At</Table.ColumnHeader>
              <Table.ColumnHeader>Updated At</Table.ColumnHeader>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {isLoading ? (
              <Table.Row>
                <Table.Cell colSpan={3} textAlign="center">
                  Loading...
                </Table.Cell>
              </Table.Row>
            ) : error ? (
              <Table.Row>
                <Table.Cell colSpan={3}>Error: {error.message}</Table.Cell>
              </Table.Row>
            ) : (
              getPagedData()?.map((task) => (
                <Table.Row key={task.id}>
                  <Table.Cell>{task.id}</Table.Cell>
                  <Table.Cell>{task.title}</Table.Cell>
                  <Table.Cell>{TaskStatus[task.status]}</Table.Cell>
                  <Table.Cell>
                    {convertUnixToDate(task.createdAt) ?? "N/A"}
                  </Table.Cell>
                  <Table.Cell>
                    {convertUnixToDate(task.updatedAt) ?? "N/A"}
                  </Table.Cell>
                </Table.Row>
              ))
            )}
          </Table.Body>
        </Table.Root>
        <Flex justify="center" mt="4">
          {filteredData && filteredData.length > 0 && !isLoading && !error && (
            <Pagination.Root
              count={filteredData!.length}
              pageSize={5}
              page={page}
              onPageChange={({ page }) => setPage(page)}
            >
              <ButtonGroup variant="ghost" size="sm" wrap="wrap">
                <Pagination.PrevTrigger asChild>
                  <IconButton>
                    <LuChevronLeft />
                  </IconButton>
                </Pagination.PrevTrigger>

                <Pagination.Items
                  render={(page) => (
                    <IconButton
                      variant={{ base: "ghost", _selected: "outline" }}
                    >
                      {page.value}
                    </IconButton>
                  )}
                />

                <Pagination.NextTrigger asChild>
                  <IconButton>
                    <LuChevronRight />
                  </IconButton>
                </Pagination.NextTrigger>
              </ButtonGroup>
            </Pagination.Root>
          )}
        </Flex>
      </Box>
    </Flex>
  );
};
