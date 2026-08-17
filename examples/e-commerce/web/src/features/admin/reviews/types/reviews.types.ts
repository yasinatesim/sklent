
export type Review = {
  id: string;
  productId: string;
  authorName: string;
  rating: number;
  comment: string;
  status: "pending" | "approved" | "rejected";
  createdAt: string;
};
