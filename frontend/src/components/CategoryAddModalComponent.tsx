import { Button, Form, Modal } from "react-bootstrap";
import Category from '../model/Categories';
import { SubmitHandler, useForm } from "react-hook-form";
import ModelAction from "../model/CategoryAction";
import Action from "../model/Action";
import { useEffect } from "react";

interface Props {
  show: boolean;
  setShow: React.Dispatch<boolean>;
  setCategoryAction: React.Dispatch<ModelAction<Category>>;
}

type CategoryInput = {
  name: string,
  description: string,
}

export default function CategoryModal({ show, setShow, setCategoryAction }: Props) {
  const { register, handleSubmit, reset, formState } = useForm<CategoryInput>();

  const handleClose = () => setShow(false);

  const addCategory: SubmitHandler<CategoryInput> = (data) => {
    setCategoryAction({ t: new Category(data.name, data.description, 0), action: Action.POST });
    handleClose();
  }

  useEffect(() => {
    if (formState.isSubmitSuccessful) reset()
  }, [reset, formState])

  return (<>
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Add category</Modal.Title>
      </Modal.Header>
      <Form onSubmit={handleSubmit(addCategory)}>
        <Modal.Body>
          <Form.Control className="mb-3" placeholder="name" {...register("name")} />
          <Form.Control className="mb-3" placeholder="description" {...register("description")} />
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>
            Close
          </Button>
          <Button variant="primary" type="submit">
            Add
          </Button>
        </Modal.Footer>
      </Form>
    </Modal>
  </>
  )
}