import { Layout } from "@/components/layout";
import { TaskForm } from "@/components/TaskForm";
import { Tooltip } from "@/components/ui/tooltip";
import type { Task } from "@/lib/api/aliases";
import { $api } from "@/lib/api/client";
import { convertUnixToDate, TaskRecurrence, TaskStatus } from "@/lib/utils";
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

const TaskDetailPage = () => {
  const [task, setTask] = useState<Task | null>(null);
  const [steps] = useState<{ title: string; description: string }[]>([
    { title: "Step 1", description: "Waiting to Start" },
    { title: "Step 2", description: "In Progress" },
    { title: "Step 3", description: "Completed" },
  ]);
  const [currentStep, setCurrentStep] = useState(1);
  const [openForm, setOpenForm] = useState(false);

  const params = useParams();
  if (!params.id) {
    return <div>Task ID is required</div>;
  }

  const { data, error, isLoading } = $api.useQuery(
    "get",
    "/tasks/{task_id}",
    { params: { path: { task_id: params.id } } },
    {
      refetchOnWindowFocus: true,
      retry: 3,
    }
  );

  useEffect(() => {
    if (!isLoading && data && data.data) {
      setTask(data.data);
      setCurrentStep(data.data.status > 2 ? 3 : data.data.status + 1);
    }
  }, [error, data, isLoading]);

  return (
    <>
      <Layout title="Task Detail">
        {task && (
          <Flex w="full" h="full" justify="space-between" mt="10">
            <Flex w="3/4" direction="column">
              <Card.Root mb="6">
                <Card.Body>
                  <Flex w="full" justify="space-between">
                    <Flex direction="column" gap="2">
                      <Card.Title>{task.title}</Card.Title>
                      <Card.Description>
                        {task.description || "No description provided."}
                      </Card.Description>
                      <Tooltip content="Task Status">
                        <Badge
                          w="fit"
                          size="md"
                          variant="outline"
                          colorPalette={
                            task.status === 0
                              ? "gray"
                              : task.status === 1
                              ? "blue"
                              : task.status === 2
                              ? "green"
                              : task.status === 3
                              ? "red"
                              : "gray"
                          }
                        >
                          {TaskStatus[task.status]}
                        </Badge>
                      </Tooltip>
                    </Flex>
                    <Flex>
                      {task.status === 0 && (
                        <Button
                          w="fit"
                          variant="outline"
                          colorPalette="blue"
                          onClick={() => setOpenForm(true)}
                        >
                          Edit Task
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
                    <TextDisplay label="Task ID" value={task.id || "N/A"} />
                    <TextDisplay
                      label="Recurrence"
                      value={TaskRecurrence[task.recurrence] || "N/A"}
                    />
                    <TextDisplay
                      label="Scheduled At"
                      value={convertUnixToDate(task.scheduledTime) || "N/A"}
                    />
                    {task.createdAt && (
                      <TextDisplay
                        label="Created At"
                        value={convertUnixToDate(task.createdAt) || "N/A"}
                      />
                    )}
                    {task.updatedAt && (
                      <TextDisplay
                        label="Updated At"
                        value={convertUnixToDate(task.updatedAt) || "N/A"}
                      />
                    )}
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
                task.status === 0
                  ? "gray"
                  : task.status === 1
                  ? "blue"
                  : task.status === 2
                  ? "green"
                  : task.status === 3
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
                {task && task.id && <TaskForm task_id={task.id} />}
              </Dialog.Body>
              <Dialog.Footer />
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </>
  );
};

export default TaskDetailPage;
