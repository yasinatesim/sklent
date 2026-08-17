export type Product = {
  id: string;
  slug: string;
  titleTr: string;
  titleEn: string;
  descriptionTr: string;
  priceCents: number;
  oldPriceCents: number;
  stock: number;
  categorySlug: string;
  badge: string;
  seller: string;
  imageUrl: string;
  published: boolean;
};

export type Category = {
  id: string;
  slug: string;
  nameTr: string;
  nameEn: string;
  icon: string;
  descTr: string;
  descEn: string;
};

export type PlaceOrderItem = {
  productId: string;
  titleTr: string;
  unitCents: number;
  quantity: number;
};
