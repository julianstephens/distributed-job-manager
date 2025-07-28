import {
  createListCollection,
  Field,
  Input,
  NumberInput,
  Select,
} from "@chakra-ui/react";
import { langs } from "@uiw/codemirror-extensions-langs";
import CodeMirror, { type ReactCodeMirrorProps } from "@uiw/react-codemirror";

export const InputFormField = ({
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

export const NumberInputFormField = ({
  name,
  defaultValue,
  onChange,
}: {
  name: string;
  defaultValue?: string;
  onChange: (value: number) => void;
}) => (
  <Field.Root>
    <Field.Label textTransform="capitalize">{name}</Field.Label>
    <NumberInput.Root
      defaultValue={defaultValue}
      onValueChange={({ value }) => {
        onChange(Number.parseInt(value));
      }}
    >
      <NumberInput.Control />
      <NumberInput.Input />
    </NumberInput.Root>
  </Field.Root>
);

export const SelectFormField = ({
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

export const CodeEditorFormField = ({
  name,
  onChange,
  ...args
}: {
  name: string;
} & ReactCodeMirrorProps) => (
  <Field.Root>
    <Field.Label textTransform="capitalize">{name}</Field.Label>
    <CodeMirror extensions={[langs.go()]} onChange={onChange} {...args} />
  </Field.Root>
);
