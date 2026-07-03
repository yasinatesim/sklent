import { Suspense } from "react";
import ResetPasswordForm from "./ResetPasswordForm";

const ResetPasswordPage = () => (
  <main className="container">
    <Suspense fallback={null}>
      <ResetPasswordForm />
    </Suspense>
  </main>
);

export default ResetPasswordPage;
