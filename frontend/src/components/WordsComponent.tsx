import { Container, Form, Row, Col, Card, Button } from "react-bootstrap";
import { useEffect, useState } from "react";
import { Trash3 } from "react-bootstrap-icons";
import { Link } from "react-router-dom";
import NavBar from "./NavBar";
import Word from '../model/Word';
import "../App.css"
import ModelAction from "../model/CategoryAction";
import Action from "../model/Action";
import { SubmitHandler, useForm } from "react-hook-form";

function Words() {
  const { register, handleSubmit, formState: { errors } } = useForm<Array<Word>>();
  const [words, setWords] = useState<Array<Word>>([]);
  const [word, setWord] = useState<ModelAction<number> | null>(null);

  const addWord = () => { setWords([...words, new Word("", "")]) }
  const deleteWord = (index: number) => {
    if (index >= 0 && index < words.length) {
      setWord({ t: index, action: Action.DELETE });
    }
  }

  const onSubmit: SubmitHandler<Array<Word>> = (data) => {
    console.log(data);
  }

  useEffect(() => {
    if (word?.action === Action.DELETE) { words.splice(word.t, 1); }
    else if (word?.action === Action.POST) { }
    setWords(words);
    setWord(null);
  }, [word])

  return (
    <>
      <NavBar />
      <Form onSubmit={handleSubmit(onSubmit)}>
        <Container className="rows">
          {words.map((word, index) => {
            return (
              <Card className="mb-3" bg="Secondary" key={index}>
                <Card.Header className="mb-1 p-2 align-middle">
                  <Row>
                    <Col md="11">{index}</Col>
                    <Col md="1"><Link to="#" onClick={() => { deleteWord(index) }}><Trash3 /></Link></Col>
                  </Row>
                </Card.Header>
                <Card.Body>
                  <Row>
                    <Form.Group as={Col} md="4">
                      <Form.Control isInvalid={errors[index]?.term != null} className="border-0 border-bottom border-primary nofocus"
                        type="text" placeholder="term" defaultValue={word.term} {...register(`${index}.term`, { required: true, maxLength: 250, minLength: 2 })} />
                      <Form.Control.Feedback type="invalid">
                        Please, provide a string with more than 2 characters and less than 250
                      </Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group as={Col} md="4">
                      <Form.Control isInvalid={errors[index]?.definition != null} className="border-0 border-bottom border-primary nofocus"
                        type="text" placeholder="definition" defaultValue={word.definition} {...register(`${index}.definition`, { required: true, maxLength: 250, minLength: 2 })} />
                      <Form.Control.Feedback type="invalid">
                        Please, provide a string with more than 2 characters and less than 250
                      </Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group as={Col} md="3">
                      <Form.Control type="file" />
                    </Form.Group>
                  </Row>
                </Card.Body>
              </Card>
            )
          })}
          <Card className="mb-3" bg="Secondary">
            <Card.Body className="mb-1 p-2 align-middle">
              <Row>
                <Col md="6">
                  {words.length}
                </Col>
                <Col md="6">
                  <Button type="button" onClick={addWord}>More</Button>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Container>
        <Container className="mt-3 text-end">
          <Button type="submit">Done</Button>
        </Container>
      </Form>
    </>
  );
}

export default Words;