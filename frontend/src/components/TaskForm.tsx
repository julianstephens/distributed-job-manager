import "@/date-picker.css";
import type { Task } from "@/lib/api/aliases";
import { $api } from "@/lib/api/client";
import { TaskRecurrence } from "@/lib/utils";
import {
  Button,
  createListCollection,
  Field,
  Flex,
  Input,
  Select,
} from "@chakra-ui/react";
import { useForm } from "@tanstack/react-form";
import { useEffect, useState } from "react";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

const InputFormField = ({
  name,
  placeholder,
  onChange,
}: {
  name: string;
  placeholder?: string;
  onChange: (value: string) => void;
}) => (
  <Field.Root>
    <Field.Label textTransform="capitalize">{name}</Field.Label>
    <Input
      borderWidth={1}
      placeholder={placeholder}
      onChange={(e) => onChange(e.currentTarget.value)}
    />
  </Field.Root>
);

const SelectFormField = ({
  name,
  placeholder,
  items,
  onChange,
}: {
  name: string;
  placeholder?: string;
  items: { value: string | number; label: string }[];
  onChange: (value: any) => void;
}) => (
  <Field.Root>
    <Field.Label textTransform="capitalize">{name}</Field.Label>
    <Select.Root
      collection={createListCollection({ items: items })}
      onValueChange={(e) => onChange(e.value[0])}
    >
      <Select.HiddenSelect />
      <Select.Control>
        <Select.Trigger>
          <Select.ValueText placeholder={placeholder} />
        </Select.Trigger>
        <Select.IndicatorGroup>
          <Select.Indicator />
        </Select.IndicatorGroup>
      </Select.Control>
      <Select.Positioner>
        <Select.Content>
          {items.map((item) => (
            <Select.Item key={item.value} item={item}>
              {item.label}
              <Select.ItemIndicator />
            </Select.Item>
          ))}
        </Select.Content>
      </Select.Positioner>
    </Select.Root>
  </Field.Root>
);

export const TaskForm = ({ task_id }: { task_id: string }) => {
  const { data, isLoading, error } = $api.useQuery("get", "/tasks/{task_id}", {
    params: { path: { task_id } },
  });
  const [task, setTask] = useState<Task | null>(null);
  const { mutate } = $api.useMutation("put", "/tasks");

  const form = useForm({
    defaultValues: {
      title: task?.title ?? "",
      description: task?.description ?? "",
      scheduledTime: task?.scheduledTime
        ? new Date(task.scheduledTime * 1000).toISOString()
        : "",
      recurrence: task?.recurrence ?? TaskRecurrence[0],
    },
    onSubmit: async ({ value: values }) => {
      const putObj: Record<string, string | number | undefined> = task ?? {};
      Object.entries(values).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== "") {
          putObj[key] = value;
        }
      });
      if (typeof putObj.scheduledTime === "string") {
        putObj.scheduledTime = Math.floor(
          new Date(putObj.scheduledTime).getTime() / 1000
        );
      }

      try {
        mutate({ body: putObj as any });
        form.reset();
      } catch (error) {
        console.error("Error updating task:", error);
      }
    },
  });

  useEffect(() => {
    if (!isLoading && data && data.data) {
      setTask(data.data);
      form.reset({
        title: data.data.title,
        description: data.data.description ?? "",
        scheduledTime: data.data.scheduledTime
          ? new Date(data.data.scheduledTime * 1000).toISOString()
          : "",
        recurrence: data.data.recurrence,
      });
    }
  }, [data, isLoading, error]);

  return (
    <>
      {task && (
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <Flex direction="column" gap="3">
            <form.Field
              name="title"
              children={(field) => (
                <InputFormField
                  name={field.name}
                  placeholder={task.title}
                  onChange={field.handleChange}
                />
              )}
            />
            <form.Field
              name="description"
              children={(field) => (
                <InputFormField
                  name={field.name}
                  placeholder={task.description}
                  onChange={field.handleChange}
                />
              )}
            />
            <form.Field
              name="recurrence"
              children={(field) => (
                <SelectFormField
                  name={field.name}
                  placeholder="Select recurrence"
                  items={Object.entries(TaskRecurrence).map(
                    ([value, label]) => ({
                      value: parseInt(value),
                      label,
                    })
                  )}
                  onChange={field.handleChange}
                />
              )}
            />
            <form.Field
              name="scheduledTime"
              children={() => (
                <Field.Root>
                  <Field.Label>Scheduled Time</Field.Label>
                  <DatePicker
                    selected={
                      form.getFieldValue("scheduledTime")
                        ? new Date(form.getFieldValue("scheduledTime"))
                        : null
                    }
                    onChange={(date) =>
                      form.setFieldValue(
                        "scheduledTime",
                        date ? date.toISOString() : ""
                      )
                    }
                    showTimeSelect
                    timeFormat="HH:mm"
                    timeIntervals={15}
                    dateFormat="MMMM d, yyyy h:mm aa"
                    placeholderText={
                      task.scheduledTime
                        ? new Date(task.scheduledTime * 1000).toLocaleString()
                        : "Select scheduled time"
                    }
                    customInput={<Input w="full" border="none" />}
                  />
                </Field.Root>
              )}
            />
          </Flex>
          <Flex w="full" justify="flex-end" mt="6">
            <form.Subscribe
              selector={(state) => [state.canSubmit, state.isSubmitting]}
              children={([canSubmit, isSubmitting]) => (
                <Button
                  ml="auto"
                  size="sm"
                  type="submit"
                  loading={isSubmitting}
                  disabled={!canSubmit}
                >
                  Submit
                </Button>
              )}
            />
          </Flex>
        </form>
      )}
    </>
  );
};
