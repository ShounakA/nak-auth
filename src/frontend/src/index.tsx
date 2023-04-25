/* @refresh reload */
import { render } from 'solid-js/web';
import { Router } from "@solidjs/router";

import './index.css';
import App from './App';
import { filter, fromEvent, map } from 'rxjs';

const root = document.getElementById('root');

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got mispelled?',
  );
}

const keyStrokes$ = fromEvent(document, 'keypress').pipe(
  filter((x) => x.type == 'keypress'),
  map((x) => x as KeyboardEvent)
);
const someTypeText = "Here is some example text to copy.";
const someTypeTextLineText = `Here.\nHere is another line :)\nYo another layer`;
render(() => 
  (<Router>
    <App />
  </Router>), 
root!);
  
