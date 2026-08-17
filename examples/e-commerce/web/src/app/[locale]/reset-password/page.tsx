import { Suspense } from "react";

import ResetPasswordForm from "@/features/web/auth/components/ResetPasswordForm";

const ResetPasswordPage = () => (
  <main className="container">
    <Suspense fallback={null}>
      <ResetPasswordForm />
    </Suspense>
  </main>
);

export default ResetPasswordPage;
