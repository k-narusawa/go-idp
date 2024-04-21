import {
  Control,
  FieldValues,
  UseControllerProps,
  useController,
  useForm,
} from "react-hook-form";
import React from "react";
import "../../styles.css";

type Props = {
  className?: string | undefined;
  placeholder?: string | undefined;
  type: string;
};

type InputItemComponentProps<T extends FieldValues> = UseControllerProps<T> &
  Props;

export const Input = <T extends FieldValues>(
  props: InputItemComponentProps<T>
) => {
  const defaultControl = useForm().control as unknown as Control<T>;
  const defaultRules = {};
  const {
    name,
    control = defaultControl,
    rules = defaultRules,
    className,
    placeholder,
  } = props;
  const { field, fieldState } = useController<T>({ name, control, rules });
  const { error } = fieldState;

  return (
    <div>
      <input
        type={props.type}
        placeholder={placeholder}
        className={`
          h-12 px-4 
          block w-full rounded-md 
          border border-gray shadow-sm
          focus:border-primary-400 focus:ring focus:blue-light focus:ring-opacity-50 
          disabled:cursor-not-allowed disabled:bg-gray-light disabled:text-gray
          focus-visible: outline-none
          ${error ? "border-red-light" : ""}
          ${className}
        `}
        {...field}
      />
      {error && (
        <span className="mt-1 text-sm text-red-500">
          {error ? error.message : ""}
        </span>
      )}
    </div>
  );
};
