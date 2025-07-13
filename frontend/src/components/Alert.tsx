import { Button, CloseButton, Dialog, Portal } from "@chakra-ui/react";

export const NegativeAlertDialog = ({
  open,
  setOpen,
  buttonLabel,
  actionCallback,
  taskId,
}: {
  open: boolean;
  setOpen: (open: boolean) => void;
  buttonLabel?: string;
  actionCallback: (taskId?: string) => void;
  taskId?: string;
}) => {
  return (
    <Dialog.Root
      role="alertdialog"
      open={open}
      onOpenChange={(e) => {
        setOpen(e.open);
      }}
    >
      <Portal>
        <Dialog.Backdrop />
        <Dialog.Positioner>
          <Dialog.Content>
            <Dialog.Header>
              <Dialog.Title>Are you sure?</Dialog.Title>
            </Dialog.Header>
            <Dialog.Body>
              <p>
                This action cannot be undone. This will permanently delete your
                task and cancel any associated work.
              </p>
            </Dialog.Body>
            <Dialog.Footer>
              <Dialog.ActionTrigger asChild>
                <Button variant="outline">Cancel</Button>
              </Dialog.ActionTrigger>
              <Button
                colorPalette="red"
                onClick={() => {
                  actionCallback(taskId);
                  setOpen(false);
                }}
              >
                {buttonLabel ?? "Delete"}
              </Button>
            </Dialog.Footer>
            <Dialog.CloseTrigger asChild>
              <CloseButton size="sm" />
            </Dialog.CloseTrigger>
          </Dialog.Content>
        </Dialog.Positioner>
      </Portal>
    </Dialog.Root>
  );
};
