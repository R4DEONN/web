const burgerEl = document.querySelector('#burger');
const headerNavEl = document.querySelector('#burgerNav');
const bodyEl = document.querySelector('#body')

burgerEl.addEventListener('click', showHeaderNav)
bodyEl.addEventListener('click', hideHeaderNav)

function showHeaderNav()
{
    headerNavEl.classList.add('appearance-burger-menu');
}

function hideHeaderNav(event) {
    if(event.target !== burgerEl)
    {
        headerNavEl.classList.remove('appearance-burger-menu');
        burgerEl.classList.remove('hidden');
    }
    else
    {
        burgerEl.classList.add('hidden');
        headerNavEl.classList.add('appearance-burger-menu');
    }
}