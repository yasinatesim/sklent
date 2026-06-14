import { Suspense } from "react";
import SuccessContent from "./SuccessContent";

const SuccessPage = () => (
  <main className="container">
    <Suspense fallback={null}>
      <SuccessContent />
    </Suspense>
  </main>
);

export default SuccessPage;
