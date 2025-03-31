// 商品的實體類
export interface CommodityListResponse {
  commodity_id: number;
  commodity_name: string;
  price_range: {
    min: number;
    max: number;
  };
  total_stock: number;
  specifications_count: number;
  picture_url?: string;
}

export interface CommodityDetailResponse {
  commodity_id: number;
  commodity_name: string;
  specification_types: SpecTypeResponse[];
  commodity_specifications: CommoditySpecResponse[];
}

export interface SpecTypeResponse {
  spec_type_id: number;
  spec_type_name: string;
  specification_values: SpecValueResponse[];
}

export interface SpecValueResponse {
  spec_value_id: number;
  spec_value: string;
}

export interface CommoditySpecResponse {
  commodity_spec_id: number;
  spec_value_1_id: number;
  spec_value_1: string;
  spec_value_2_id?: number | null; // 可選，對應 *uint
  spec_value_2?: string | null;   // 可選，對應 *string
  stock: number;
  price: number;
  picture_url: string;
}
