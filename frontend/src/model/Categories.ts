export default class Category {
  name: string;
  description: string;
  amount: number;

  constructor(name: string, description: string, amount: number) {
    this.name = name;
    this.description = description;
    this.amount = amount;
  }
}