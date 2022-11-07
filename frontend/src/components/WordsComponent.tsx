import { Container, Form, Row, Col, Card, Button } from "react-bootstrap";
import NavBar from "./NavBar";
import "../App.css"
import { useState } from "react";
import { Trash3 } from "react-bootstrap-icons";
import { Link } from "react-router-dom";

function Words() {
  const [inputForms, setInputForms] = useState<number>(1);

  const addInputForm = () => { setInputForms(inputForms+1) }
  // const deleteInputForm = (index: number) => { if (index > 0 && index <= inputForms) { delete inputForms[index]; setInputForms(inputForms) } }
  const deleteInputForm = (index: number) => {};

  return (
    <>
      <NavBar />
      <Container className="rows">
        {inputForms.map((v, index) => {
          return (
            <Card className="mb-3" bg="Secondary" key={index}>
              <Card.Header className="mb-1 p-2 align-middle">
                <Row>
                  <Col className="col-11">{index}</Col>
                  <Col className="col-1"><Link to="#" onClick={deleteInputForm(index)}><Trash3 /></Link></Col>
                </Row>
              </Card.Header>
              <Card.Body>
                <Row>
                  <Col>
                    <Form.Control className="border-0 border-bottom border-primary nofocus" type="text" placeholder="term" />
                  </Col>
                  <Col>
                    <Form.Control className="border-0 border-bottom border-primary nofocus" type="text" placeholder="definition" />
                  </Col>
                  <Col className="col-3">
                    <Form.Control type="file" />
                  </Col>
                </Row>
              </Card.Body>
            </Card>
          )
        })}
      </Container>
      <Container className="mt-3 text-center">
        <Button type="button" onClick={addInputForm}>More</Button>
      </Container>
    </>
  );
}

export default Words;