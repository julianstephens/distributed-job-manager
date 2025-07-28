import { JobForm } from "@/components/JobForm";
import { Layout } from "@/components/layout";
import { Tooltip } from "@/components/ui/tooltip";
import { useJob } from "@/lib/api/hooks";
import type { Job } from "@/lib/types";
import {
  Badge,
  Button,
  Card,
  CloseButton,
  Dialog,
  Flex,
  Portal,
  Separator,
  Steps,
  Text,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useParams } from "react-router";

const TextDisplay = ({ label, value }: { label: string; value: string }) => (
  <Flex gap="2">
    <Flex direction="column" w="1/4">
      <Text>{label}</Text>
    </Flex>
    <Flex direction="column" w="1/4">
      <Text>:</Text>
    </Flex>
    <Flex direction="column">
      <Text color="gray.400">{value}</Text>
    </Flex>
  </Flex>
);

const JobDetailPage = () => {
  const [job, setJob] = useState<Job | null>(null);
  const [steps] = useState<{ title: string; description: string }[]>([
    { title: "Step 1", description: "Waiting to Start" },
    { title: "Step 2", description: "In Progress" },
    { title: "Step 3", description: "Completed" },
  ]);
  const [currentStep, setCurrentStep] = useState(1);
  const [openForm, setOpenForm] = useState(false);

  const params = useParams();
  if (!params.id) {
    return <div>Job ID is required</div>;
  }

  const { data, error, isLoading } = useJob(params.id);

  useEffect(() => {
    if (!isLoading && data && data) {
      setJob(data);
      setCurrentStep(1);
    }
  }, [error, data, isLoading]);

  return (
    <>
      <Layout title="Job Detail">
        {job && (
          <Flex w="full" h="full" justify="space-between" mt="10">
            <Flex w="3/4" direction="column">
              <Card.Root mb="6">
                <Card.Body>
                  <Flex w="full" justify="space-between">
                    <Flex direction="column" gap="2">
                      <Card.Title>{job.job_name}</Card.Title>
                      <Tooltip content="Job Status">
                        <Badge
                          w="fit"
                          size="md"
                          variant="outline"
                          colorPalette={
                            job.status === "pending"
                              ? "gray"
                              : job.status === "in-progress"
                              ? "blue"
                              : job.status === "completed"
                              ? "green"
                              : job.status === "failed"
                              ? "red"
                              : "gray"
                          }
                        >
                          {job.status}
                        </Badge>
                      </Tooltip>
                    </Flex>
                    <Flex>
                      {job.status === "pending" && (
                        <Button
                          w="fit"
                          variant="outline"
                          colorPalette="blue"
                          onClick={() => setOpenForm(true)}
                        >
                          Edit Job
                        </Button>
                      )}
                    </Flex>
                  </Flex>
                </Card.Body>
              </Card.Root>
              <Flex w="full" gap="2">
                <Card.Root w="1/3">
                  <Card.Body>
                    <Card.Title mb="4">Task Details</Card.Title>
                    <TextDisplay label="Task ID" value={job.job_id} />
                    <TextDisplay label="Recurrence" value={job.frequency} />
                    <TextDisplay
                      label="Next Execution Time"
                      value={job.execution_time}
                    />
                    {/* {job.createdAt && (
                      <TextDisplay
                        label="Created At"
                        value={convertUnixToDate(job.createdAt) || "N/A"}
                      />
                    )}
                    {job.updatedAt && (
                      <TextDisplay
                        label="Updated At"
                        value={convertUnixToDate(job.updatedAt) || "N/A"}
                      />
                    )} */}
                  </Card.Body>
                </Card.Root>
                <Card.Root w="2/3">
                  <Card.Body>
                    <Card.Title mb="4">Run History</Card.Title>
                    <Text mx="auto" my="auto">
                      Nothing to display
                    </Text>
                  </Card.Body>
                </Card.Root>
              </Flex>
            </Flex>
            <Separator orientation="vertical" mx="10" />
            <Steps.Root
              orientation="vertical"
              height={"400px"}
              count={steps.length}
              defaultStep={currentStep}
              variant="subtle"
              colorPalette={
                job.status === "pending"
                  ? "gray"
                  : job.status === "in-progress"
                  ? "blue"
                  : job.status === "completed"
                  ? "green"
                  : job.status === "failed"
                  ? "red"
                  : "gray"
              }
              w="1/4"
            >
              <Steps.List>
                {steps.map((step, index) => (
                  <Steps.Item key={index} index={index} title={step.title}>
                    <Steps.Indicator />
                    <Steps.Title>{step.description}</Steps.Title>
                    <Steps.Separator />
                  </Steps.Item>
                ))}
              </Steps.List>
            </Steps.Root>
          </Flex>
        )}
        {error && <div>Error loading task: {error.message}</div>}
      </Layout>
      <Dialog.Root
        open={openForm}
        onOpenChange={(details) => setOpenForm(details.open)}
      >
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.CloseTrigger asChild>
                <CloseButton size="sm" />
              </Dialog.CloseTrigger>
              <Dialog.Header>
                <Dialog.Title>Edit Task</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                {job && job.job_id && <JobForm job_id={job.job_id} />}
              </Dialog.Body>
              <Dialog.Footer />
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </>
  );
};

export default JobDetailPage;
