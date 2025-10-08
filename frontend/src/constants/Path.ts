const SETTING_BASE = "/setting";

const PATH = {
  // Auth & General
  HOME: "/",
  SIGN_IN: "/sign-in",
  FIRST_CHANGE_PASSWORD: "/first-change-password",

  // User
  PROFILE: "/profile",

  CHECKIN: "/attendant",

  // HR
  CONTRACT: "/contract",
  DECISION: "/decision",
  INSURANCE: "/insurance",
  APPROVE: "/approve",
  DOCUMENT_DETAIL: "/documentDetail",
  DETAIL_DECISION: "/detail_decision",
  SHIFT_MANAGEMENT: "/shift_management",
  CREATESHIFT: "/shift_create",
  SHIFTDETAILPAGE: "/ShiftDetailPage",
  ROTA: "/rota",
  // HR Checkin
  WORKSHIFT: "/workshift",
  ATTENDANCE_MANAGEMENT: "/attendance-management",
  APPROVE_ATTENDANCE: "/approve-attendant",
  TIMESHEET: "/timesheet",
  ATTENDANCE_REPORT: "/attendant-report",
  ATTENDANT_HISTORY: "/attendant-history",
  ATTENDANT_HISTORY_BY_DATE: "/attendant-history/by-date",
  ATTENDANT_HISTORY_EMPLOYEES: "/attendant-history/employees",
  WORK_SCHEDULE: "/setup-work-schedule",
  WORK_SCHEDULE_REGISTER: "/setup-work-schedule-register",
  // HR Setting
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
  REGISTER_WORKSHIFT: "/register-workshift",
  DOCUMENTDETAIL: "/documentDetail",

  // Manager
  PROFILE_MANAGER: "/profile-manager",
  ATTENDANT_HISTORY_MANAGER: "/attendant-history-manager",
  WORKSHIFT_EMPLOYEE: "/workshift-employee",
};

export default PATH;
