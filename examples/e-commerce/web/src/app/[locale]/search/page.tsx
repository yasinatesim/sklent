import { Suspense } from "react";

import SearchContent from "@/features/web/catalog/components/SearchContent";

import styles from "@/features/web/catalog/styles/search.module.scss";

const SearchPage = () => (
  <main className={`container ${styles.page}`}>
    <Suspense fallback={null}>
      <SearchContent />
    </Suspense>
  </main>
);

export default SearchPage;
