import "@/date-picker.css";
import {
  createListCollection,
  Field,
  Input,
  NumberInput,
  Select,
  Textarea,
} from "@chakra-ui/react";
import { langs } from "@uiw/codemirror-extensions-langs";
import CodeMirror, { type ReactCodeMirrorProps } from "@uiw/react-codemirror";
import type { PropsWithChildren } from "react";
import DatePicker from "react-datepicker";

type FormFieldProps = {
  name: string;
  error?: string;
};

const FormField = ({
  name,
  error,
  children,
}: FormFieldProps & PropsWithChildren) => {
  return (
    <Field.Root invalid={!!error}>
      <Field.Label textTransform="capitalize">{name}</Field.Label>
      {children}
      {error && <Field.ErrorText>{error}</Field.ErrorText>}
    </Field.Root>
  );
};

export const InputFormField = ({
  onChange,
  placeholder,
  ...args
}: FormFieldProps & {
  onChange: (value: string) => void;
  placeholder?: string;
}) => (
  <FormField {...args}>
    <Input
      borderWidth={1}
      placeholder={placeholder}
      onChange={(e) => onChange(e.currentTarget.value)}
    />
  </FormField>
);

export const TextAreaFormField = ({
  placeholder,
  onChange,
  ...args
}: FormFieldProps & {
  placeholder?: string;
  onChange: (value: string) => void;
}) => (
  <FormField {...args}>
    <Textarea
      borderWidth={1}
      placeholder={placeholder}
      onChange={(e) => onChange(e.currentTarget.value)}
    />
  </FormField>
);

export const NumberInputFormField = ({
  defaultValue,
  onChange,
  ...args
}: FormFieldProps & {
  defaultValue?: string;
  onChange: (value: number) => void;
}) => (
  <FormField {...args}>
    <NumberInput.Root
      defaultValue={defaultValue}
      onValueChange={({ value }) => {
        onChange(Number.parseInt(value));
      }}
    >
      <NumberInput.Control />
      <NumberInput.Input />
    </NumberInput.Root>
  </FormField>
);

export const SelectFormField = ({
  placeholder,
  defaultValue,
  items,
  onChange,
  ...args
}: FormFieldProps & {
  placeholder?: string;
  defaultValue?: string;
  items: { value: string | number; label: string }[];
  onChange: (value: any) => void;
}) => (
  <FormField {...args}>
    <Select.Root
      collection={createListCollection({ items: items })}
      onValueChange={(e) => onChange(e.value[0])}
      defaultValue={defaultValue ? [defaultValue] : undefined}
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
  </FormField>
);

export const CodeEditorFormField = ({
  name,
  onChange,
  error,
  ...args
}: FormFieldProps & ReactCodeMirrorProps) => (
  <FormField name={name} error={error}>
    <CodeMirror extensions={[langs.go()]} onChange={onChange} {...args} />
  </FormField>
);

export const DatePickerFormField = ({
  placeholder,
  onChange,
  selected,
  error,
  name,
  onBlur,
}: FormFieldProps & {
  placeholder?: string;
  selected: Date;
  onChange: (date: Date) => void;
  onBlur?: any;
}) => (
  <FormField error={error} name={name}>
    <DatePicker
      selected={selected}
      onChange={onChange as any}
      onBlur={onBlur}
      showTimeSelect
      timeFormat="HH:mm"
      timeIntervals={15}
      dateFormat="MMMM d, yyyy h:mm aa"
      placeholderText={placeholder}
      customInput={<Input w="full" border="none" />}
      className={error && "border-red"}
    />
  </FormField>
);
