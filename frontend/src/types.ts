export interface ChildrenProps {
  children: React.ReactNode;
}

export const TaskStatus = {
  0: "Pending",
  1: "In Progress",
  2: "Completed",
  3: "Failed",
  4: "Cancelled",
} as const;
