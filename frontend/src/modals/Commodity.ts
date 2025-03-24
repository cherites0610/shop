class Commodity {
  id: number;
  name: string;
  price: number;
  specification: { [key: string]: string[] };
  stock: number;

  constructor(id:number, name: string, price: number, specification: { [key: string]: string[] }, stock: number) {
    this.id = id;
    this.name = name;
    this.price = price;
    this.specification = specification;
    this.stock = stock;
  }
}

export default Commodity;
