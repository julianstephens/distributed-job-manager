import { Layout } from "@/components/layout";
import { TaskTable } from "@/components/TaskTable";

const LandingPage = () => {
  return (
    <Layout title="Task List">
      <TaskTable />
    </Layout>
  );
};

export default LandingPage;
