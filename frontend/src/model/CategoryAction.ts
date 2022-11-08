import Action from "./Action";

type ModelAction<T> = {
  t: T;
  action: Action;
}

export default ModelAction;