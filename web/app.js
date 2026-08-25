// --------------------------------------------------
// Кнопки
// --------------------------------------------------

document
    .querySelectorAll("button[data-analytics-name]")
    .forEach(button => {

        button.addEventListener("click", () => {

            analytics.track(
                "click",
                button,
                {
                    text: button.innerText
                }
            );

        });

    });


// --------------------------------------------------
// Текстовые поля
// --------------------------------------------------

document
    .querySelectorAll(
        "input[type='text'], input[type='email']"
    )
    .forEach(input => {

        input.addEventListener("input", () => {

            analytics.track(
                "input",
                input,
                {
                    value: input.value
                }
            );

        });

    });


// --------------------------------------------------
// Checkbox
// --------------------------------------------------

document
    .querySelectorAll("input[type='checkbox']")
    .forEach(checkbox => {

        checkbox.addEventListener("change", () => {

            analytics.track(
                "change",
                checkbox,
                {
                    checked: checkbox.checked
                }
            );

        });

    });


// --------------------------------------------------
// Radio buttons
// --------------------------------------------------

document
    .querySelectorAll("input[type='radio']")
    .forEach(radio => {

        radio.addEventListener("change", () => {

            analytics.track(
                "change",
                radio,
                {
                    value: radio.value
                }
            );

        });

    });


// --------------------------------------------------
// Select
// --------------------------------------------------

document
    .querySelectorAll("select[data-analytics-name]")
    .forEach(select => {

        select.addEventListener("change", () => {

            analytics.track(
                "change",
                select,
                {
                    value: select.value
                }
            );

        });

    });


// --------------------------------------------------
// Slider
// --------------------------------------------------

const volume = document.querySelector("#volume");
const volumeValue = document.querySelector("#volume-value");

volume.addEventListener("input", () => {

    volumeValue.innerText = volume.value;

    analytics.track(
        "slider",
        volume,
        {
            value: Number(volume.value)
        }
    );

});