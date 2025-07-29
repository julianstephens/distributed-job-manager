import { toaster } from "@/components/ui/toaster";
import "@/date-picker.css";
import { useAuthToken, useJob } from "@/lib/api/hooks";
import { createJob } from "@/lib/api/queries";
import type { Job, JobRequest } from "@/lib/types";
import { getKeyByValue, JobFrequency, queryClient } from "@/lib/utils";
import { Button, Field, Flex, Input } from "@chakra-ui/react";
import { useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import { tokyoNight } from "@uiw/codemirror-theme-tokyo-night";
import { useEffect, useState } from "react";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import {
  CodeEditorFormField,
  InputFormField,
  NumberInputFormField,
  SelectFormField,
} from "./form";

export const JobForm = ({
  job_id,
  closeForm,
}: {
  job_id?: string;
  closeForm?: VoidFunction;
}) => {
  const token = useAuthToken();
  const { data, isLoading, error } = useJob(job_id);
  const [job, setJob] = useState<Job | null>(null);
  const mutation = useMutation({
    mutationFn: createJob,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["jobs"] });
      if (job_id) {
        queryClient.invalidateQueries({ queryKey: ["jobs", job_id] });
      }
    },
  });

  const form = useForm({
    defaultValues: {
      job_name: job?.job_name ?? "",
      frequency: job?.frequency ?? JobFrequency.Once,
      payload: job?.payload ?? "\n".repeat(7),
      max_retries: job?.max_retries ?? 3,
      execution_time: job?.execution_time ?? new Date().toISOString(),
    },
    onSubmit: async ({ value: values }) => {
      const req: JobRequest = {
        job_name: values.job_name,
        frequency: values.frequency,
        payload: "```go\n" + values.payload + "\n```",
        max_retries: values.max_retries,
        execution_time: values.execution_time,
      };

      try {
        const job = await mutation.mutateAsync({ token, req });
        toaster.success({
          title: `Job ${job.job_name} created`,
          description: `Scheduled job ${job.job_name} to run ${getKeyByValue(
            JobFrequency,
            job.frequency as any
          )?.toLowerCase()} starting ${job.execution_time}`,
        });
        form.reset();
        if (closeForm) closeForm();
      } catch (error) {
        toaster.error({ title: "Unable to create job", description: error });
      }
    },
  });

  useEffect(() => {
    if (!isLoading && data) {
      setJob(data);
      form.reset({
        job_name: data.job_name,
        frequency: data.frequency,
        payload: data.payload,
        max_retries: data.max_retries,
        execution_time: data.execution_time,
      });
    }
  }, [data, isLoading, error]);

  return (
    <>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}
      >
        <Flex direction="column" gap="3">
          <form.Field
            name="job_name"
            children={(field) => (
              <InputFormField
                name="name"
                placeholder={job?.job_name}
                onChange={field.handleChange}
              />
            )}
          />
          <form.Field
            name="frequency"
            children={(field) => (
              <SelectFormField
                name={field.name}
                placeholder="Select frequency"
                items={Object.entries(JobFrequency).map(([label, value]) => ({
                  value,
                  label,
                }))}
                onChange={field.handleChange}
              />
            )}
          />
          <form.Field
            name="max_retries"
            children={(field) => (
              <NumberInputFormField
                name={field.name.split("_").join(" ")}
                defaultValue={`${job?.max_retries ?? 3}`}
                onChange={field.handleChange}
              />
            )}
          />
          <form.Field
            name="execution_time"
            children={(field) => (
              <Field.Root>
                <Field.Label textTransform="capitalize">
                  {field.name.split("_").join(" ")}
                </Field.Label>
                <DatePicker
                  selected={
                    form.getFieldValue(field.name)
                      ? //@ts-ignore TS2769
                        new Date(form.getFieldValue(field.name))
                      : new Date()
                  }
                  onChange={(date) =>
                    form.setFieldValue(
                      "execution_time",
                      date ? date.toISOString() : ""
                    )
                  }
                  showTimeSelect
                  timeFormat="HH:mm"
                  timeIntervals={15}
                  dateFormat="MMMM d, yyyy h:mm aa"
                  placeholderText={
                    job?.execution_time
                      ? job.execution_time
                      : "Select execution time"
                  }
                  customInput={<Input w="full" border="none" />}
                />
              </Field.Root>
            )}
          />
          <form.Field
            name="payload"
            children={(field) => (
              <CodeEditorFormField
                name={field.name}
                onChange={field.handleChange}
                height="200px"
                width="100%"
                editable={true}
                theme={tokyoNight}
                placeholder={form.getFieldValue(field.name) as string}
                indentWithTab={true}
                basicSetup={{
                  lineNumbers: true,
                  highlightActiveLine: true,
                  bracketMatching: true,
                  closeBrackets: true,
                  autocompletion: true,
                }}
              />
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
    </>
  );
};
