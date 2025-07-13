import { Layout } from "@/components/layout";
import type { Task } from "@/lib/api/aliases";
import { $api } from "@/lib/api/client";
import { TaskStatus } from "@/lib/utils";
import {
  Badge,
  Button,
  Flex,
  Heading,
  Separator,
  Steps,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useParams } from "react-router";

const TaskDetailPage = () => {
  const [task, setTask] = useState<Task | null>(null);
  const [steps] = useState<{ title: string; description: string }[]>([
    { title: "Step 1", description: "Waiting to Start" },
    { title: "Step 2", description: "In Progress" },
    { title: "Step 3", description: "Completed" },
  ]);
  const [currentStep, setCurrentStep] = useState(1);

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
      console.log("Fetched task data:", data);
      setTask(data.data);
      setCurrentStep(data.data.status > 2 ? 3 : data.data.status + 1);
    }
  }, [error, data, isLoading]);

  return (
    <Layout title="Task Detail">
      {task && (
        <Flex justify="space-between" mt="10">
          <Flex w="3/4" h="fit" align="center" justify="space-between">
            <Flex direction="column">
              <Heading>{task.title}</Heading>
              {task?.description && <p>{task.description}</p>}
            </Flex>
            <Flex gap="4">
              <Badge
                size="lg"
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
              {task.status === 0 && (
                <Button variant="outline" colorPalette="blue">
                  Edit Task
                </Button>
              )}
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
  );
};

export default TaskDetailPage;
