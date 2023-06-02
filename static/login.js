const form = document.forms[0];
const inputEmailEl = document.querySelector('#inputEmail');
const inputPasswordEl = document.querySelector('#inputPassword');
const emailErrorEl = document.querySelector('#emailError');
const passwordErrorEl = document.querySelector('#passwordError')
const formErrorEl = document.querySelector('#mainError');
const formErrorTextEl = document.querySelector('#mainErrorText');
const passwordEyeEl = document.querySelector('#passwordEye');
const formTitleEl = document.querySelector('#formTitle');

initListeners();

function initListeners()
{
  inputPasswordEl.addEventListener('input', changeStyle);
  inputEmailEl.addEventListener('input', changeStyle);
  inputPasswordEl.addEventListener('click', removeError);
  inputEmailEl.addEventListener('click', removeError);
  passwordEyeEl.addEventListener('click', changeEyeStyle);
}

function changeStyle(event)
{
  const el = event.target;
  if (el.value !== '')
  {
    el.classList.add('form-row__input_filled');
  }
  else
  {
    el.classList.remove('form-row__input_filled');
  }
}

function removeError(event)
{
  const el = event.target;
  formTitleEl.style.marginBottom = '50px';
  if (el.classList.contains('form-row__input_error'))
  {
    el.classList.remove('form-row__input_error');
  }
  const lastChild = el.parentNode.lastChild.previousSibling;
  if (lastChild.nodeName === 'SPAN')
  {
    lastChild.classList.add('hidden');
  }
  const nextParentSibling = el.parentNode.nextSibling.nextSibling;
  if (nextParentSibling.nodeName === 'SPAN')
  {
    nextParentSibling.classList.add('hidden');
  }
  if (formErrorEl.classList.contains('appearance'))
  {
    formErrorEl.classList.remove('appearance');
  }
}

function changeEyeStyle(event)
{
  const el = event.target;
  if (inputPasswordEl.type === "password")
  {
    el.src = '/static/img/eye-off.svg';
    inputPasswordEl.type = "text";
  }
  else
  {
    el.src = '/static/img/eye.svg';
    inputPasswordEl.type = "password";
  }
}

function validateQueryParams()
{
  let isErrors = false;
  if (inputEmailEl.value === '')
  {
    isErrors = true;
    emailErrorEl.innerHTML = 'Email is required.';
    emailErrorEl.classList.remove('hidden');
  }
  else if (!inputEmailEl.value.match(/^[-a-z0-9!#$%&'*+/=?^_`{|}~]+(\.[-a-z0-9!#$%&'*+/=?^_`{|}~]+)*@([a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?\.)+[a-z]+/))
  {
    isErrors = true;
    emailErrorEl.innerHTML = 'Incorrect email format. Correct format is ****@**.***';
    emailErrorEl.classList.remove('hidden');
  }

  if (inputPasswordEl.value === '')
  {
    isErrors = true;
    passwordErrorEl.innerHTML = 'Password is required.';
    passwordErrorEl.classList.remove('hidden');
  }

  return isErrors;
}


form.onsubmit = async e =>
{
  e.preventDefault();
  let isErrors= validateQueryParams();
  if (isErrors)
  {
    formErrorTextEl.innerHTML = 'A-Ah! Check all fields.';
    if (!formErrorEl.classList.contains('appearance'))
    {
      formErrorEl.classList.add('appearance');
    }
    formTitleEl.style.marginBottom = '40px';
    return;
  }

  formTitleEl.style.marginBottom = '50px';
  const props = {};
  for (let element of form.elements)
  {
    if (element.type === 'submit') continue;
    props[element.name] = element.value;
  }
  
  const json = JSON.stringify(props, null, '\t');
  console.log(json);

  let response = await fetch('/api/login', {
        method: 'POST',
        headers: {
        'Content-Type': 'application/json;charset=utf-8'
        },
        body: json
    });

  if (response.ok)
  {
    window.location = '/admin';
  }
  else
  {
    if (!formErrorEl.classList.contains('appearance'))
    {
      formErrorEl.classList.add('appearance');
    }
    formErrorTextEl.innerHTML = 'Email or password is incorrect.';
    formTitleEl.style.marginBottom = '40px';
    inputPasswordEl.classList.add('form-row__input_error');
    inputEmailEl.classList.add('form-row__input_error');
  }
}