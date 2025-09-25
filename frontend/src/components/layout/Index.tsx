import "./index.scss";
import React from "react";
import Header from "./components/Header";
import Content from "./components/Content";
const Home: React.FC = () => {
    return (
        <div className="app-layout">
            <Header />
            <Content />
        </div>
    );
};

export default Home;
