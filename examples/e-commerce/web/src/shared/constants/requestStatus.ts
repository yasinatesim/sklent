export const REQUEST_STATUS = {
  IDLE: "idle",
  LOADING: "loading",
  SUCCESS: "success",
  ERROR: "error",
} as const;

export type RequestStatus = (typeof REQUEST_STATUS)[keyof typeof REQUEST_STATUS];
