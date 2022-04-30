import { io } from "socket.io-client";

//TODO: give correct url
const socket = io('http://localhost:8000/socket.io/');

export default socket;