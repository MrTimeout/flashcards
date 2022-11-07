import { useEffect, useState } from "react";
import { Button, Container, Table } from "react-bootstrap";
import Category from "../model/categories";
import CategoryModal from "./CategoryAddModalComponent";
import service from '../service/category';
import { PencilSquare, Trash3 } from "react-bootstrap-icons";
import { Link } from "react-router-dom";
import NavBar from "./NavBar";
import CategoryAction from "../model/categoryAction";
import Action from "../model/action";

function Categories() {
  // const [categories, setCategories] = useState(new Array<Category>())
  // const [categories, setCategories] = useState<Category[]>([])
  const [categories, setCategories] = useState(new Array<Category>());
  const [categoryAction, setCategoryAction] = useState<CategoryAction | null>(null);
  const [showModal, setShowModal] = useState(false);
  const showCategoryModal = () => setShowModal(true);

  useEffect(() => {
    if (categoryAction != null) {
      if (categoryAction.action === Action.POST) service.post(categoryAction.c);
      else if (categoryAction.action === Action.DELETE) service.remove(categoryAction.c.name)
      setCategoryAction(null);
    }
    service.getAll().then((response: any) => {
      setCategories(response.data);
    }).catch((e: Error) => {
      console.log(e);
    });
  }, [categoryAction])


  return (
    <>
      <NavBar />
      <Container>
        <Table striped bordered hover variant="light">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Amount</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {
              categories.map((c: Category, i: number) => {
                return (
                  <tr key={i}>
                    <td>{c.name}</td>
                    <td>{c.description}</td>
                    <td className="text-center">{c.amount ?? 0}</td>
                    <td className="text-center">
                      <Link to={`/categories/${c.name}`} className="me-1"><PencilSquare /></Link>
                      <Link to="#" onClick={() => setCategoryAction({ c: c, action: Action.DELETE })}><Trash3 /></Link>
                    </td>
                  </tr>
                )
              })
            }
          </tbody>
        </Table>
        <Button variant="primary" type="button" onClick={showCategoryModal}>Add Category</Button>
      </Container>
      <CategoryModal show={showModal} setShow={setShowModal} setCategoryAction={setCategoryAction} />
    </>
  )
}

export default Categories;