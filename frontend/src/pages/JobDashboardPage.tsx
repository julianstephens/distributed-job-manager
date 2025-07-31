import { JobForm } from "@/components/JobForm";
import { JobTable } from "@/components/JobTable";
import { Layout } from "@/components/layout";
import { useAuth0 } from "@auth0/auth0-react";
import { Button, CloseButton, Dialog, Portal } from "@chakra-ui/react";
import { useState } from "react";

const JobDashboardPage = () => {
  const { user } = useAuth0();
  const [openForm, setOpenForm] = useState(false);

  const createButton = () => (
    <Button
      onClick={() => {
        setOpenForm(true);
      }}
    >
      New Job
    </Button>
  );

  return (
    <>
      <Layout
        title={`Hello, ${user?.given_name ?? "User"}`}
        actionButton={createButton()}
      >
        <JobTable />
      </Layout>
      <Dialog.Root
        open={openForm}
        onOpenChange={(details) => setOpenForm(details.open)}
      >
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.CloseTrigger asChild>
                <CloseButton size="sm" />
              </Dialog.CloseTrigger>
              <Dialog.Header>
                <Dialog.Title>New Job</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                {
                  <JobForm
                    closeForm={() => {
                      setOpenForm(false);
                    }}
                  />
                }
              </Dialog.Body>
              <Dialog.Footer />
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
    </>
  );
};
export default JobDashboardPage;
