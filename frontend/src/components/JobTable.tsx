import { Tooltip } from "@/components/ui/tooltip";
import { useJobs } from "@/lib/api/hooks";
import type { Job } from "@/lib/types";
import {
  displayDate,
  getJobStatusColor,
  JobStatus,
  TABLE_PAGE_SIZE,
} from "@/lib/utils";
import { useAuth0 } from "@auth0/auth0-react";
import {
  Badge,
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
import { FaArrowDown, FaArrowUp, FaBan, FaEye, FaX } from "react-icons/fa6";
import { LuChevronLeft, LuChevronRight } from "react-icons/lu";
import { NegativeAlertDialog } from "./Alert";

const ColumnWithSortArrow = ({
  field,
  sortKey,
  sortBy,
  sortDirection,
}: {
  field: string;
  sortKey: string;
  sortBy: string;
  sortDirection: string;
}) => (
  <Flex justify="center" align="center" gap="2">
    {field}
    <Icon display={sortBy === sortKey ? "inline-block" : "none"}>
      {sortDirection === "asc" ? <FaArrowUp /> : <FaArrowDown />}
    </Icon>
  </Flex>
);

export const JobTable = () => {
  const { user } = useAuth0();
  const { data, error, isLoading } = useJobs(user?.sub);

  const sortFilters = createListCollection({
    items: [
      { label: "Name", value: "name" },
      { label: "Status", value: "status" },
      { label: "Execution Time", value: "executionTime" },
      { label: "Created At", value: "createdAt" },
      { label: "Last Updated", value: "lastUpdated" },
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
  const [filteredData, setFilteredData] = useState<Job[] | null>(null);
  const [page, setPage] = useState(1);
  const [openCancelDialog, setOpenCancelDialog] = useState(false);
  const [jobContext, setJobContext] = useState<string | undefined>(undefined);
  const [availableStatuses, setAvailableStatuses] = useState<string[]>([]);

  const updateSort = (sortBy: string | null) => {
    if (!filteredData) return;
    if (!sortBy) {
      setFilteredData(data || null);
      return;
    }
    const sortedData = [...filteredData].sort((a, b) => {
      if (sortBy === "name") {
        return a.job_name < b.job_name ? -1 : a.job_name === b.job_name ? 0 : 1;
      } else if (sortBy === "executionTime") {
        return (
          new Date(a.execution_time).getTime() -
          new Date(b.execution_time).getTime()
        );
      } else if (sortBy === "createdAt") {
        return (
          new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
        );
      } else if (sortBy === "lastUpdated") {
        return (
          new Date(a.updated_at).getTime() - new Date(b.updated_at).getTime()
        );
      } else if (sortBy === "status") {
        return a.status < b.status ? -1 : a.status === b.status ? 0 : 1;
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
      setFilteredData(data || null);
      setStatusFilter(null);
      return;
    }

    setFilteredData(data?.filter((job) => job.status === status) || null);
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
    if (!isLoading && !error && data) {
      setFilteredData(data);
      setAvailableStatuses(Array.from(new Set(data.map((job) => job.status))));
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
                items={Object.values(JobStatus).map((status) => ({
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
                <Table.ColumnHeader>
                  <ColumnWithSortArrow
                    field="Name"
                    sortKey="name"
                    sortBy={sortBy as string}
                    sortDirection={sortDirection as string}
                  />
                </Table.ColumnHeader>
                <Table.ColumnHeader>Description</Table.ColumnHeader>
                <Table.ColumnHeader>
                  <ColumnWithSortArrow
                    field="Status"
                    sortKey="status"
                    sortBy={sortBy as string}
                    sortDirection={sortDirection as string}
                  />
                </Table.ColumnHeader>
                <Table.ColumnHeader>
                  <ColumnWithSortArrow
                    field="Execution Time"
                    sortKey="executionTime"
                    sortBy={sortBy as string}
                    sortDirection={sortDirection as string}
                  />
                </Table.ColumnHeader>
                <Table.ColumnHeader>
                  <ColumnWithSortArrow
                    field="Created At"
                    sortKey="createdAt"
                    sortBy={sortBy as string}
                    sortDirection={sortDirection as string}
                  />
                </Table.ColumnHeader>
                <Table.ColumnHeader>
                  <ColumnWithSortArrow
                    field="Last Updated"
                    sortKey="lastUpdated"
                    sortBy={sortBy as string}
                    sortDirection={sortDirection as string}
                  />
                </Table.ColumnHeader>
                <Table.ColumnHeader color="gray">Actions</Table.ColumnHeader>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {isLoading ? (
                <Table.Row>
                  <Table.Cell colSpan={8} textAlign="center">
                    Loading...
                  </Table.Cell>
                </Table.Row>
              ) : error ? (
                <Table.Row>
                  <Table.Cell colSpan={8} textAlign="center">
                    Error: {error.message}
                  </Table.Cell>
                </Table.Row>
              ) : !filteredData || filteredData.length === 0 ? (
                <Table.Row>
                  <Table.Cell colSpan={8} textAlign="center">
                    No jobs to display
                  </Table.Cell>
                </Table.Row>
              ) : (
                getPagedData()?.map((job) => (
                  <Table.Row key={job.job_id}>
                    <Table.Cell>{job.job_name}</Table.Cell>
                    <Table.Cell>{job.job_description}</Table.Cell>
                    <Table.Cell>
                      <Badge
                        padding="2"
                        colorPalette={getJobStatusColor(job.status as any)}
                        fontSize="sm"
                        borderRadius="lg"
                        display="flex"
                        justifyContent="center"
                        alignItems="center"
                      >
                        {job.status}
                      </Badge>
                    </Table.Cell>
                    <Table.Cell>{displayDate(job.execution_time)}</Table.Cell>
                    <Table.Cell>{displayDate(job.created_at)}</Table.Cell>
                    <Table.Cell>{displayDate(job.updated_at)}</Table.Cell>
                    <Table.Cell>
                      <ButtonGroup variant="outline" size="xs">
                        <Link
                          href={`/jobs/${job.job_id}`}
                          style={{ textDecoration: "none" }}
                          _hover={{ fill: "blue" }}
                        >
                          <IconButton aria-label="View Job Details">
                            <FaEye className="icon-hover-purple" />
                          </IconButton>
                        </Link>
                        {!["completed", "failed", "cancelled"].includes(
                          job.status
                        ) && (
                          <IconButton
                            aria-label="Cancel Job"
                            rotate="90deg"
                            onClick={() => {
                              setOpenCancelDialog(true);
                              setJobContext(job.job_id);
                            }}
                          >
                            <FaBan className="icon-hover-red" />
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
        buttonLabel="Cancel Job"
        jobId={jobContext}
        actionCallback={(jobId?: string) => {
          // TODO: Implement job cancellation API call
          console.log("Cancel Job", jobId);
        }}
      />
    </>
  );
};
