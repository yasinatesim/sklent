"use client";

import { useCallback, useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import AdminProductForm from "@/components/admin/AdminProductForm";
import { fetchProducts, fetchCategories, formatTRY, type Product, type Category } from "@/lib/api";
import styles from "../admin.module.scss";

const AdminProductsPage = () => {
  const t = useTranslations("admin");
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);

  const loadProducts = useCallback(async () => {
    try {
      setProducts(await fetchProducts());
    } catch {
      setProducts([]);
    }
  }, []);

  useEffect(() => {
    const loadCategories = async () => {
      try {
        setCategories(await fetchCategories());
      } catch {
        setCategories([]);
      }
    };
    loadCategories();
    loadProducts();
  }, [loadProducts]);

  return (
    <div>
      <h2>{t("products")}</h2>
      <AdminProductForm categories={categories} onCreated={loadProducts} />
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("product")}</th>
            <th>{t("category")}</th>
            <th>{t("price")}</th>
            <th>{t("stock")}</th>
          </tr>
        </thead>
        <tbody>
          {products.map((p) => (
            <tr key={p.id}>
              <td>{p.titleTr}</td>
              <td>{p.categorySlug || "—"}</td>
              <td>{formatTRY(p.priceCents)}</td>
              <td>{p.stock}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminProductsPage;
