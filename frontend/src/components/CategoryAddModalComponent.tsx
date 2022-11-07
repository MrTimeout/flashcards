import { Button, Form, Modal } from "react-bootstrap";
import Category from '../model/categories';
import { SubmitHandler, useForm } from "react-hook-form";
import CategoryAction from "../model/categoryAction";
import Action from "../model/action";
import { useEffect } from "react";

interface Props {
  show: boolean;
  setShow: React.Dispatch<boolean>;
  setCategoryAction: React.Dispatch<CategoryAction>;
}

type CategoryInput = {
  name: string,
  description: string,
}

export default function CategoryModal({ show, setShow, setCategoryAction }: Props) {
  const { register, handleSubmit, reset, formState } = useForm<CategoryInput>();

  const handleClose = () => setShow(false);

  const addCategory: SubmitHandler<CategoryInput> = (data) => {
    setCategoryAction({ c: new Category(data.name, data.description, 0), action: Action.POST });
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