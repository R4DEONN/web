const TITLE_ARRAY = [document.querySelector('.article-preview__title'), document.querySelector('.card-preview__title')];
const SUBTITLE_ARRAY = [document.querySelector('.article-preview__subtitle'), document.querySelector('.card-preview__subtitle')];
const AVATAR_ARRAY = [document.querySelector('.form-row__avatar-placeholder'), document.querySelector('.card-preview__author-image')];
const MAIN_IMAGE_ARRAY = [document.querySelector('.upload__main-image'), document.querySelector('.article-preview__image')];
const PREVIEW_IMAGE_ARRAY = [document.querySelector('.upload__preview-image'), document.querySelector('.card-preview__image')];


const logOutButton = document.querySelector('#logOutButton');
logOutButton.addEventListener('click', () => { window.location = '/login' })

const textInputElements = document.querySelectorAll('.form-row__input');
for (let el of textInputElements)
{
  el.addEventListener('change', ChangeStyle);
}
function ChangeStyle(event)
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


const form = document.forms[0];
const alertMessage = document.querySelector('#alertMessage');
const successMessage = document.querySelector('#successMessage');

form.onsubmit = async e =>
{
  e.preventDefault();
  let errors = ValidateQueryParams(form.elements);
  if (Object.keys(errors).length !== 0)
  {
    if (alertMessage.classList.contains('hidden'))
    {
      alertMessage.classList.remove('hidden');
    }
    if (!successMessage.classList.contains('hidden'))
    {
      successMessage.classList.add('hidden');
    }
    const json = JSON.stringify(errors, null, '\t');
    console.log(json);
    return;
  }

  if (!alertMessage.classList.contains('hidden'))
  {
    alertMessage.classList.add('hidden');
  }
  if (successMessage.classList.contains('hidden'))
  {
    successMessage.classList.remove('hidden');
  }
  const props = {};
  for (let element of form.elements)
  {
    if (element.type == 'submit') continue;
    props[element.name] = element.value;
  }
  const json = JSON.stringify(props, null, '\t');
  console.log(json);

  fetch('/createPost', {
        method: 'POST',
        headers: {
        'Content-Type': 'application/json;charset=utf-8'
        },
        body: json
    });
}

function ValidateQueryParams(query)
{
  let errors = {};
  for (let element of query)
  {
    if (element.value === '' && element.type !== 'submit')
    {
      errors[element.name] = 'Поле ' + element.name + ' не должно быть пустым';
    }
  }
  return errors;
}

const titleEl = document.querySelector('#inputTitle');
titleEl.addEventListener('change', previewTitle);
function previewTitle(event)
{
  for (let element of TITLE_ARRAY)
  {
    element.innerHTML = event.target.value;
  }
}

const subtitleEl = document.querySelector('#inputSub');
subtitleEl.addEventListener('change', previewSubtitle);
function previewSubtitle(event)
{
  for (let element of SUBTITLE_ARRAY)
  {
    element.innerHTML = event.target.value;
  }
}

const dateEl = document.querySelector('#inputDate');
dateEl.addEventListener('change', previewDate);

function previewDate(event)
{
  let element = document.querySelector('.card-preview__date');
  element.innerHTML = event.target.value;
}

const authorNameEl = document.querySelector('#inputAuthorName');
authorNameEl.addEventListener('change', previewAuthorName);

function previewAuthorName(event)
{
  let element = document.querySelector('.card-preview__author-name');
  if (event.target.value)
  {
    element.innerHTML = event.target.value;
  }
  else
  {
    element.innerHTML = 'Enter author name';
  }
}

const authorImageEl = document.querySelector('#inputAuthorImage');
authorImageEl.addEventListener('change', previewAuthorImage);
const uploadAvatarButtonText = document.querySelector('#uploadAvatarText');
const avatarCameraImg = document.querySelector('#avatarCamera');
const removeAvatarButton = document.querySelector('#removeAvatarButton');
removeAvatarButton.addEventListener('click', deleteAvatar);

function previewAuthorImage(event)
{
  const reader = new FileReader();
  reader.onloadend = function ()
  {
    if (reader.result === '')
    {
      return;
    }

    avatarCameraImg.classList.remove('hidden');
    uploadAvatarButtonText.innerHTML = 'Upload New';
    removeAvatarButton.classList.remove('hidden');
    for (let image of AVATAR_ARRAY)
    {
      image.style.background = "url(" + reader.result + ")";
      image.style.backgroundSize = "cover";
    }

  }

  if (event.target.files[0])
  {
    reader.readAsDataURL(event.target.files[0]);
  }
  else
  {
    deleteAvatar();
  }
}

function deleteAvatar()
{
  removeAvatarButton.classList.add('hidden');
  avatarCameraImg.classList.add('hidden');
  authorImageEl.value = '';
  uploadAvatarButtonText.innerHTML = 'Upload';
  AVATAR_ARRAY[0].style.backgroundImage = '';
  AVATAR_ARRAY[1].style.background = "#F7F7F7";
}


const mainImageEl = document.querySelector('#inputMainImage');
mainImageEl.addEventListener('change', previewMainImage);

const mainImageController = document.querySelector('#mainImageController');
const removeMainImageButton = document.querySelector('#removeMainImageButton');
removeMainImageButton.addEventListener('click', deleteMainImage);
const mainImageRemark = document.querySelector('#mainImageRemark');


function previewMainImage(event)
{
  let reader = new FileReader();
  reader.onloadend = function ()
  {
    if (reader.result === '')
    {
      return;
    }
    if (mainImageController.classList.contains('hidden'))
    {
      mainImageRemark.classList.add('hidden');
      mainImageController.classList.remove('hidden');
    }
    for (let image of MAIN_IMAGE_ARRAY)
    {
      image.style.background = "url(" + reader.result + ")"
      image.style.backgroundSize = "cover";
      image.classList.add('upload__main-image_uploaded')
    }
  }

  if (event.target.files[0])
  {
    reader.readAsDataURL(event.target.files[0]);
  }
  else
  {
    deleteMainImage();
  }
}

function deleteMainImage()
{
  mainImageController.classList.add('hidden');
  mainImageRemark.classList.remove('hidden');
  mainImageEl.value = '';
  MAIN_IMAGE_ARRAY[0].style.backgroundImage = '';
  MAIN_IMAGE_ARRAY[0].classList.remove('upload__main-image_uploaded');
  MAIN_IMAGE_ARRAY[1].style.background = "#F7F7F7";
  MAIN_IMAGE_ARRAY[1].classList.remove('upload__main-image_uploaded');
}


const previewImageEl = document.querySelector('#inputPreviewImage');
previewImageEl.addEventListener('change', previewPreviewImage);


const previewImageController = document.querySelector('#previewImageController');
const removePreviewImageButton = document.querySelector('#removePreviewImageButton');
removePreviewImageButton.addEventListener('click', deletePreviewImage);
const previewImageRemark = document.querySelector('#previewImageRemark');

function previewPreviewImage(event)
{
  let reader = new FileReader();
  reader.onloadend = function ()
  {
    if (reader.result === '')
    {
      return;
    }
    if (previewImageController.classList.contains('hidden'))
    {
      previewImageRemark.classList.add('hidden');
      previewImageController.classList.remove('hidden');
    }
    for (let image of PREVIEW_IMAGE_ARRAY)
    {
      image.style.background = "url(" + reader.result + ")"
      image.style.backgroundSize = "cover";
      image.classList.add('upload__preview-image_uploaded')
    }
  }

  if (event.target.files[0])
  {
    reader.readAsDataURL(event.target.files[0]);
  }
  else
  {
    deletePreviewImage();
  }
}

function deletePreviewImage()
{
  previewImageController.classList.add('hidden');
  previewImageRemark.classList.remove('hidden');
  previewImageEl.value = '';
  PREVIEW_IMAGE_ARRAY[0].style.backgroundImage = '';
  PREVIEW_IMAGE_ARRAY[0].classList.remove('upload__preview-image_uploaded');
  PREVIEW_IMAGE_ARRAY[1].style.background = "#F7F7F7";
  PREVIEW_IMAGE_ARRAY[1].classList.remove('upload__preview-image_uploaded');
}