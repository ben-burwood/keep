import { createApp } from 'vue'
import App from './App.vue'
import './style.css';

export const SERVER_URL = `${window.location.protocol}//${window.location.hostname}:8080`;

createApp(App).mount('#app')
