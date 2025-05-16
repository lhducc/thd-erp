const createPath = () => {
  const PROFILE = "/profile";
  const CONTRACT = "/contract";
  const DECISION = "/decision";
  return {
    HOME: "/",
    SIGN_IN: "/sign-in",
    FIRSRT_CHANGE_PASSWORD: "/first-change-password",
    PROFILE,
    CONTRACT,
    DECISION,
    OFFICE: "/office",
    DEPARTMENT: "/department",
    POSITION: "/position",
    JOB_TITLE: "/job-title",
    HIERARCHY_LEVEL: "/hierarchy-level",
  };
};

const PATH = createPath();

export default PATH;
