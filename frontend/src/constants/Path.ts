const SETTING_BASE = "/setting";

const PATH = {
  HOME: "/",
  SIGN_IN: "/sign-in",
  FIRST_CHANGE_PASSWORD: "/first-change-password",
  PROFILE: "/profile",
  CONTRACT: "/contract",
  DECISION: "/decision",
  INSURANCE: "/insurance",

  SETTING: SETTING_BASE,
  SETTING_OFFICE: `${SETTING_BASE}/office`,
  SETTING_DEPARTMENT: `${SETTING_BASE}/department`,
  SETTING_POSITION: `${SETTING_BASE}/position`,
  SETTING_JOB_TITLE: `${SETTING_BASE}/job-title`,
  SETTING_HIERARCHY_LEVEL: `${SETTING_BASE}/hierarchy-level`,
  SETTING_CONTRACT: `${SETTING_BASE}/contract`,
  SETTING_ALLOWANCE: `${SETTING_BASE}/allowance`,
  SETTING_EMPLOYEE_DOCUMENT: `${SETTING_BASE}/employee-document`,
};

export default PATH;
