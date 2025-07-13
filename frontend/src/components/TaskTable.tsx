import { Tooltip } from "@/components/ui/tooltip";
import type { Task } from "@/lib/api/aliases";
import { $api } from "@/lib/api/client";
import { convertUnixToDate, TABLE_PAGE_SIZE, TaskStatus } from "@/lib/utils";
import {
  Box,
  ButtonGroup,
  createListCollection,
  Flex,
  Icon,
  IconButton,
  Link,
  Pagination,
  Portal,
  SegmentGroup,
  Select,
  Table,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { FaBan, FaEye, FaX } from "react-icons/fa6";
import { LuChevronLeft, LuChevronRight } from "react-icons/lu";
import { NegativeAlertDialog } from "./Alert";

export const TaskTable = () => {
  const { data, error, isLoading } = $api.useQuery("get", "/tasks", undefined, {
    refetchOnWindowFocus: true,
  });

  const sortFilters = createListCollection({
    items: [
      { label: "Title", value: "title" },
      { label: "Status", value: "status" },
      { label: "Created At", value: "createdAt" },
      { label: "Update At", value: "updatedAt" },
    ],
  });
  const sortDirections = createListCollection({
    items: [
      { label: "Ascending", value: "asc" },
      { label: "Descending", value: "desc" },
    ],
  });
  const [sortBy, setSortBy] = useState<string | null>(null);
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");
  const [statusFilter, setStatusFilter] = useState<string | null>(null);
  const [filteredData, setFilteredData] = useState<Task[] | null>(null);
  const [page, setPage] = useState(1);
  const [openCancelDialog, setOpenCancelDialog] = useState(false);
  const [taskContext, setTaskContext] = useState<string | undefined>(undefined);
  const [availableStatuses, setAvailableStatuses] = useState<string[]>([]);

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
    setSortBy(sortBy);
  };

  const updateSortDirection = (direction: "asc" | "desc") => {
    if (!filteredData) return;
    const sortedData = [...filteredData].reverse();
    setFilteredData(sortedData);
    setSortDirection(direction);
  };

  const updateStatusFilter = (status: string | null) => {
    if (!status) {
      setFilteredData(data?.data || null);
      setStatusFilter(null);
      return;
    }

    setFilteredData(
      data?.data?.filter((task) => TaskStatus[task.status] === status) || null
    );
    setStatusFilter(status);
    setPage(1);
  };

  const getPagedData = () => {
    if (!filteredData) return null;
    const startIndex = (page - 1) * TABLE_PAGE_SIZE;
    const endIndex = startIndex + TABLE_PAGE_SIZE;
    return filteredData.slice(startIndex, endIndex);
  };

  useEffect(() => {
    if (!isLoading && !error && data && data.data) {
      setFilteredData(data.data);
      setAvailableStatuses(
        Array.from(new Set(data.data.map((task) => TaskStatus[task.status])))
      );
    }
  }, [data, isLoading, error]);

  return (
    <>
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
                  <Select.ClearTrigger onClick={() => setSortBy(null)} />
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
                disabled={isLoading || sortBy === null}
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
                items={Object.values(TaskStatus).map((status) => ({
                  label: status,
                  value: status,
                  disabled: !availableStatuses.includes(status) || isLoading,
                }))}
                cursor="pointer"
                _disabled={{ cursor: "not-allowed" }}
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
                <Table.ColumnHeader color="gray">Actions</Table.ColumnHeader>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {isLoading ? (
                <Table.Row>
                  <Table.Cell colSpan={6} textAlign="center">
                    Loading...
                  </Table.Cell>
                </Table.Row>
              ) : error ? (
                <Table.Row>
                  <Table.Cell colSpan={6}>Error: {error.message}</Table.Cell>
                </Table.Row>
              ) : !filteredData || filteredData.length === 0 ? (
                <Table.Row>
                  <Table.Cell colSpan={6}>No tasks to display</Table.Cell>
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
                    <Table.Cell>
                      <ButtonGroup variant="outline" size="xs">
                        <Link
                          href={`/tasks/${task.id}`}
                          style={{ textDecoration: "none" }}
                        >
                          <IconButton aria-label="View Task Details">
                            <FaEye />
                          </IconButton>
                        </Link>
                        {task.status < 2 && (
                          <IconButton
                            aria-label="Cancel Task"
                            rotate="90deg"
                            onClick={() => {
                              setOpenCancelDialog(true);
                              setTaskContext(task.id);
                            }}
                          >
                            <FaBan />
                          </IconButton>
                        )}
                      </ButtonGroup>
                    </Table.Cell>
                  </Table.Row>
                ))
              )}
            </Table.Body>
          </Table.Root>
          <Flex id="paginationContainer" justify="center" mt="4">
            {filteredData &&
              filteredData.length > 0 &&
              !isLoading &&
              !error && (
                <Pagination.Root
                  count={filteredData!.length}
                  pageSize={TABLE_PAGE_SIZE}
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
      <NegativeAlertDialog
        open={openCancelDialog}
        setOpen={setOpenCancelDialog}
        buttonLabel="Cancel Task"
        taskId={taskContext}
        actionCallback={(taskId?: string) => {
          // TODO: Implement task cancellation API call
          console.log("Cancel Task", taskId);
        }}
      />
    </>
  );
};
