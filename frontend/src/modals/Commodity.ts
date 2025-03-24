// 規格的實體類
export interface Specification {
  spec_value: string[];  // 規格值，例如 ["綠色", "X"]
  stock: number;         // 庫存數量
  price: number;         // 價格
}

// 商品的實體類
export default interface Commodity {
  id: number;           // 商品 ID
  name: string;         // 商品名稱
  spec: { [key: string]: string[] };    // 規格鍵，例如 ["顏色", "大小"]
  specifications: Specification[]; // 規格列表
}