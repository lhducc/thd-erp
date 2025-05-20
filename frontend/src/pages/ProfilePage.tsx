import ButtonCreateProfile from "@/components/ButtonCreateProfile";
import DataTable from "@/components/DataTable";

const ProfilePage = () => {
  const columns = [];
  return (
    <DataTable
      columns={[]}
      data={[]}
      title="Hồ sơ nhân viên"
      buttonCreate={<ButtonCreateProfile />}
    />
  );
};

export default ProfilePage;
