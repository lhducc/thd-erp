const SETTING_BASE = "/setting";

const PATH = {
  // Auth & General
  HOME: "/",
  SIGN_IN: "/login",
  FIRST_CHANGE_PASSWORD: "/first-change-password",

  // User
  PROFILE: "/profile",

  // Pages
  CONTRACT: "/contract",
  DECISION: "/decision",
  INSURANCE: "/insurance",
  // DETAIL_INFO: "/detail_infor",
  DOCUMENT_DETAIL: "/documentDetail",
  DETAIL_DECISION: "/detail_decision",
  SHIFT_MANAGEMENT: "/shift_management",
  CREATESHIFT: "/shift_create",
  SHIFTDETAILPAGE: "/ShiftDetailPage",
  CHECKIN: "/checkin",

  SETTING: SETTING_BASE,
  SETTING_OFFICE: `${SETTING_BASE}/office`,
  SETTING_DEPARTMENT: `${SETTING_BASE}/department`,
  SETTING_DECISION: `${SETTING_BASE}/decision`,
  SETTING_POSITION: `${SETTING_BASE}/position`,
  SETTING_JOB_TITLE: `${SETTING_BASE}/job-title`,
  SETTING_HIERARCHY_LEVEL: `${SETTING_BASE}/hierarchy-level`,
  SETTING_CONTRACT: `${SETTING_BASE}/contract`,
  SETTING_ALLOWANCE: `${SETTING_BASE}/allowance`,
  SETTING_EMPLOYEE_DOCUMENT: `${SETTING_BASE}/employee-document`,

  // Other Paths
  DETAIL_INFO: "/detail_infor",
  DOCUMENTDETAIL: "/documentDetail",
};

export default PATH;
