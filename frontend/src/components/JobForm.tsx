import { toaster } from "@/components/ui/toaster";
import { useAuthToken, useJob } from "@/lib/api/hooks";
import { createJob, updateJob } from "@/lib/api/queries";
import type { Job, JobCreateRequest, JobUpdateRequest } from "@/lib/types";
import {
  formatPayload,
  getKeyByValue,
  JobFrequency,
  queryClient,
} from "@/lib/utils";
import { Button, Flex } from "@chakra-ui/react";
import { useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import { tokyoNight } from "@uiw/codemirror-theme-tokyo-night";
import axios from "axios";
import { useEffect, useState } from "react";
import "react-datepicker/dist/react-datepicker.css";
import {
  CodeEditorFormField,
  DatePickerFormField,
  InputFormField,
  NumberInputFormField,
  SelectFormField,
  TextAreaFormField,
} from "./form";

const cleanHtml = (unsanitized?: string) => {
  if (!unsanitized) return unsanitized;
  const div = document.createElement("div");
  div.innerHTML = unsanitized;
  return div.innerText;
};

const isValidDate = (value: string) => {
  const now = new Date();
  const fiveMinutesFromNow = new Date(now.getTime() + 5 * 60 * 1000);
  const date = new Date(value);

  if (date.getTime() < fiveMinutesFromNow.getTime()) {
    return false;
  }
  return true;
};

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
  const createMutation = useMutation({
    mutationFn: createJob,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["jobs"] });
    },
  });
  const updateMutation = useMutation({
    mutationFn: updateJob,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["jobs"] });
      queryClient.invalidateQueries({ queryKey: ["jobs", job_id] });
    },
  });

  const createMutate = async (values: Record<string, any>) => {
    const req: JobCreateRequest = {
      job_name: values.job_name,
      job_description: values.job_description,
      frequency: values.frequency,
      payload: formatPayload(values.payload),
      max_retries: values.max_retries,
      execution_time: values.execution_time,
    };

    try {
      const job = await createMutation.mutateAsync({ token, req });
      toaster.success({
        title: `Job ${job.job_name} created`,
        description: `${job.job_name} scheduled to run ${getKeyByValue(
          JobFrequency,
          job.frequency as any
        )?.toLowerCase()} starting ${job.execution_time}`,
      });
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        toaster.error({
          title: "Unable to create job",
          description: error.message,
        });
        return;
      }

      toaster.error({
        title: "Unable to create job",
        description: "Something went wrong",
      });
    }
  };

  const updateMutate = async (
    id: string,
    values: Record<string, any>,
    original: Job | null
  ) => {
    if (!original) {
      toaster.error({
        title: "Unable to update job",
        description: "Something went wrong.",
      });
      return;
    }

    const req: JobUpdateRequest = {
      ...(values.job_name != original.job_name && {
        job_name: values.job_name,
      }),
      ...(values.job_description != original.job_description && {
        job_description: values.job_description,
      }),
      ...(values.frequency != original.frequency && {
        frequency: values.frequency,
      }),
      ...(values.max_retries != original.max_retries && {
        max_retries: values.max_retries,
      }),
      ...(values.execution_time != original.execution_time && {
        execution_time: values.execution_time,
      }),
      payload: formatPayload(values.payload),
    };

    try {
      const job = await updateMutation.mutateAsync({ token, req, id });
      toaster.success({
        title: `Job ${job.job_name} updated`,
        description: `${job.job_name} scheduled to run ${getKeyByValue(
          JobFrequency,
          job.frequency as any
        )?.toLowerCase()} starting ${job.execution_time}`,
      });
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        toaster.error({
          title: "Unable to update job",
          description: error.message,
        });
        return;
      }

      toaster.error({
        title: "Unable to update job",
        description: "Something went wrong",
      });
    }
  };

  const form = useForm({
    defaultValues: {
      job_name: job?.job_name ?? "",
      job_description: job?.job_description ?? "",
      frequency: job?.frequency ?? JobFrequency.Once,
      payload: cleanHtml(job?.payload) ?? "\n".repeat(7),
      max_retries: job?.max_retries ?? 3,
      execution_time: job?.execution_time ?? new Date().toISOString(),
    },
    onSubmit: async ({ value: values }) => {
      if (job_id) {
        await updateMutate(job_id, values, job);
      } else {
        await createMutate(values);
      }
      if (closeForm) closeForm();
      form.reset();
    },
  });

  useEffect(() => {
    if (!isLoading && data) {
      setJob(data);
      form.reset({
        job_name: data.job_name,
        job_description: data.job_description,
        frequency: data.frequency,
        payload: cleanHtml(data.payload) ?? "\n".repeat(7),
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
          void form.handleSubmit();
        }}
      >
        <Flex direction="column" gap="3">
          <form.Field
            name="job_name"
            validators={{
              onChange: ({ value }) =>
                !value || value.length < 3
                  ? "Jobs must have a name at least 3 characters long"
                  : undefined,
            }}
            children={(field) => (
              <InputFormField
                name="name"
                placeholder={job?.job_name}
                onChange={field.handleChange}
                error={field.state.meta.errors.join(",")}
              />
            )}
          />
          <form.Field
            name="job_description"
            children={(field) => (
              <TextAreaFormField
                name="description"
                placeholder={job?.job_description}
                onChange={field.handleChange}
                error={field.state.meta.errors.join(",")}
              />
            )}
          />
          <form.Field
            name="frequency"
            children={(field) => (
              <SelectFormField
                name={field.name}
                placeholder="Select frequency"
                defaultValue={
                  form.getFieldValue(field.name) as string | undefined
                }
                error={field.state.meta.errors.join(",")}
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
                error={field.state.meta.errors.join(",")}
              />
            )}
          />
          <form.Field
            name="execution_time"
            validators={{
              onBlur: ({ value }) => {
                if (!value) return "Execution datetime must be set";
                if (!isValidDate(value))
                  return "Execution datetime must be at least 5 minutes in the future";
                return undefined;
              },
            }}
            children={(field) => (
              <DatePickerFormField
                name={field.name.split("_").join(" ")}
                selected={
                  form.getFieldValue(field.name)
                    ? //@ts-ignore TS2769
                      new Date(form.getFieldValue(field.name))
                    : new Date()
                }
                onBlur={field.handleBlur}
                onChange={(date) =>
                  form.setFieldValue(
                    "execution_time",
                    date ? date.toISOString() : ""
                  )
                }
                placeholder={
                  job?.execution_time
                    ? job.execution_time
                    : "Select execution time"
                }
                error={field.state.meta.errors.join(",")}
              />
            )}
          />
          <form.Field
            name="payload"
            validators={{
              onChange: ({ value }) => {
                if (!value || !value.replaceAll("\n", "").replaceAll(" ", ""))
                  return "Jobs must include a payload";
                return undefined;
              },
            }}
            children={(field) => (
              <CodeEditorFormField
                name={field.name}
                error={
                  field.state.meta.errors.length > 0
                    ? field.state.meta.errors[0]
                    : ""
                }
                onChange={field.handleChange}
                height="200px"
                width="100%"
                editable={true}
                theme={tokyoNight}
                value={form.getFieldValue(field.name) as string}
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
            selector={(state) => [
              state.canSubmit,
              state.isSubmitting,
              state.isValid,
            ]}
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
