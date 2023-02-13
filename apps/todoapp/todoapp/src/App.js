import React, { useState, useEffect } from "react";
import './App.css';
import axios from "axios";

const TodoApp = () => {
  const [todos, setTodos] = useState([]);
  const [newTodo, setNewTodo] = useState("");
  const [selectedTodo, setSelectedTodo] = useState(null);

  useEffect(() => {
    axios.get("/todos").then((res) => {
      setTodos(res.data);
    });
  }, []);

  const handleSubmit = (event) => {
    event.preventDefault();
    if (!selectedTodo) {
      axios.post("/todos", { text: newTodo }).then((res) => {
        setTodos([...todos, res.data]);
        setNewTodo("");
      });
    } else {
      axios.put(`/todos/${selectedTodo.id}`, { text: newTodo }).then((res) => {
        setTodos(
          todos.map((todo) => (todo.id === selectedTodo.id ? res.data : todo))
        );
        setSelectedTodo(null);
        setNewTodo("");
      });
    }
  };

  const handleEdit = (todo) => {
    setSelectedTodo(todo);
    setNewTodo(todo.text);
  };

  const handleDelete = (id) => {
    axios.delete(`/todos/${id}`).then((res) => {
      setTodos(todos.filter((todo) => todo.id !== id));
    });
  };

  return (
    <div class="container">
      <h1>To-Do App</h1>
      <form onSubmit={handleSubmit}>
        <input
	  class="edit-input"
          type="text"
          value={newTodo}
          onChange={(event) => setNewTodo(event.target.value)}
        />
        <button class="button complete-button" type="submit">{selectedTodo ? "Update" : "Add"}</button>
      </form>
      <ul class="todo-list">
        {todos.map((todo) => (
          <li class="todo-item" key={todo.id}>
            <button class="button complete-button" onClick={() => handleEdit(todo)}>Edit</button>
            <button class="button delete-button" onClick={() => handleDelete(todo.id)}>Delete</button>
	    <div class="todo-text">
              {todo.text}
	    </div>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TodoApp;

