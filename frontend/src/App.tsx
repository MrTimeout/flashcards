import Categories from './components/CategoriesComponent';
import 'bootstrap/dist/css/bootstrap.min.css';
import { Route, Routes } from 'react-router-dom';
import Words from './components/WordsComponent';

function App() {
  return (
    <Routes>
      <Route path="/categories" element={<Categories />} />
      <Route path="/categories/:name" element={<Words />} />
    </Routes>
  );
}

export default App;
