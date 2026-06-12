import { Suspense } from "react";
import SearchContent from "./SearchContent";
import styles from "./page.module.scss";

const SearchPage = () => (
  <main className={`container ${styles.page}`}>
    <Suspense fallback={null}>
      <SearchContent />
    </Suspense>
  </main>
);

export default SearchPage;
