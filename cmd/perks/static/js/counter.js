// Bucket Counter
class Bucket{
    store=localStorage
    name="cmsBucket"
    state={
        items:{},
    }

    constructor(name){
        if (name)
            this.name=name;
        const data = this.store.getItem(this.name)
        const state = JSON.parse(data);
        if (state)
            this.state=state;
    }

    getItem = (key) => {
        return this.state.items[key];
    }

    setItem = (key, count) =>{
        this.state.items[key]=count;
        this.store.setItem(this.name, JSON.stringify(this.state));
    }

    rmItem = (key) =>{
        delete(this.state.items, key);
        this.store.setItem(this.name, JSON.stringify(this.state));
    }

}


class Counter{
    count=0;
    input=false
    constructor(el, name){
        this.el = el;
        this.bucket = new Bucket(name);
        this.key=el.dataset["key"];
        const count = this.bucket.getItem(this.key);
        if (count){
            this.count=count;
        }
        this.render();
        this.setter(this.count);
    }

    minus(){
        this.count--;
        if (this.count<0)
            this.count=0;
        this.setter(this.count);
    }
    plus(e){
        this.count++;
        this.setter(this.count);
    }
    setter(value){
        if (value==0)
            this.bucket.rmItem(this.key);
        this.bucket.setItem(this.key, value)
        this.input.setAttribute("value", value);
    }

    render(){
        const btn1 = document.createElement("button");
        btn1.appendChild(document.createTextNode("-"));
        btn1.onclick=this.minus.bind(this);
        this.el.appendChild(btn1);
        this.input = document.createElement("input");
        this.el.appendChild(this.input);

        const btn2 = document.createElement("button");
        btn2.appendChild(document.createTextNode("+"));
        btn2.onclick=this.plus.bind(this);
        this.el.appendChild(btn2);
    }

}

const els = document.getElementsByClassName("bucket_counter");
if (els.length)
    Array.prototype.forEach.call(els, el=>{
        new Counter(el);
    })


