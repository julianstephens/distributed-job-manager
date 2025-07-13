import type { ChildrenProps } from "@/lib/types";
import { Box, Heading } from "@chakra-ui/react";

export const Layout = ({
  children,
  title,
}: ChildrenProps & { title: string }) => {
  return (
    <Box m={10}>
      <Heading size="2xl" mb="5">
        {title}
      </Heading>
      {children}
    </Box>
  );
};
