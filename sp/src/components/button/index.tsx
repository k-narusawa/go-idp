export const Button = ({
  children,
  variant = "primary",
  size = "default",
  disabled,
  ...props
}: {
  children: React.ReactNode;
  variant?: string;
  size?: string;
  disabled: boolean;
  [key: string]: any;
}) => {
  let baseStyle = "w-full py-2 px-4 font-normal rounded-lg shadow-md";
  if (size === "small") {
    baseStyle = "w-full py-1 px-2 font-small rounded-md shadow-sm";
  } else if (size === "large") {
    baseStyle = "w-full py-3 px-6 font-large rounded-xl shadow-lg";
  }

  let variantStyle = "";
  if (disabled) {
    variantStyle = "text-gray bg-gray-light cursor-not-allowed";
  } else if (variant === "primary") {
    variantStyle = "text-white bg-blue hover:bg-blue-dark";
  } else if (variant === "secondary") {
    variantStyle = "text-gray-dark bg-gray-light hover:bg-gray";
  } else if (variant === "alternative") {
    variantStyle = "text-black bg-white hover:bg-gray-light";
  } else if (variant === "danger") {
    variantStyle = "text-white-dark bg-red hover:bg-red-dark";
  }

  return (
    <button
      className={`${baseStyle} ${variantStyle}`}
      {...props}
      disabled={disabled}
    >
      {children}
    </button>
  );
};
