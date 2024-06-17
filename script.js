// Load the CSV file
const csvFile = './output.csv';
let csvData = [];

// Function to load the CSV file
function loadCSV() {
    fetch(csvFile)
        .then(response => response.text())
        .then(data => {
            const rows = data.split('\n');
            csvData = rows.map(row => row.split(','));
            // Initialize quiz questions array
            quizQuestions = [];
            // Loop through each row in the CSV file
            csvData.forEach((row, index) => {
                // Create a new quiz question object
                const question = {
                    id: row[0],
                    type: 'selection', // or 'signe', or 'write'
                    options: []
                };
                // Add options to the question object based on the row data
                switch (row.length) {
                    case 3: // selection type
                        question.type = 'selection';
                        question.options = [row[1], row[2], row[3], row[0]];
                        break;
                    case 2: // signe type
                        question.type = 'signe';
                        question.options = [row[0], row[1]];
                        break;
                    default: // write type
                        question.type = 'write';
                        break;
                }
                quizQuestions.push(question);
            });
            // Shuffle the questions array to randomize the order
            quizQuestions = shuffle(quizQuestions);
        });
}

// Function to shuffle an array
function shuffle(array) {
    for (let i = array.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
}

// Function to display the next question in the quiz
function displayNextQuestion() {
    if (quizQuestions.length === 0) {
        console.log('No more questions!');
        return;
    }
    const question = quizQuestions.shift();
    const questionHtml = '';
    switch (question.type) {
        case 'selection':
            questionHtml += `
        <h2>${question.id}</h2>
        <img src="${question.id}.gif" alt="${question.id}">
        <p>Choose the correct sign:</p>
        <ul>
          ${question.options.map(option => `<li>${option}</li>`).join('')}
        </ul>
      `;
            break;
        case 'signe':
            questionHtml += `
        <h2>${question.id}</h2>
        <p>Which GIF corresponds to this sign:</p>
        <ul>
          ${csvData.map(row => `<li><img src="${row[0].gif}" alt="${row[0]}></li>`).join('')}
        </ul>
      `;
            break;
        case 'write':
            questionHtml += `
        <h2>${question.id}</h2>
        <p>Write the sign for this GIF:</p>
        <input type="text" id="answer">
      `;
            break;
    }
    document.getElementById('quiz-container').innerHTML = questionHtml;
}

loadCSV();
displayNextQuestion();

// Event listener for submitting answers
document.getElementById('quiz-container').addEventListener('click', (event) => {
    if (event.target.tagName === 'LI') {
        const selectedOption = event.target.textContent;
        // Check if the answer is correct and update the quiz accordingly
        if (selectedOption === question.options[0]) { // adjust this based on your CSV data structure
            console.log('Correct!');
        } else {
            console.log('Incorrect!');
        }
        displayNextQuestion();
    } else if (event.target.tagName === 'INPUT') { // for write type questions
        const answerText = event.target.value;
        // Check if the answer is correct and update the quiz accordingly
        if (answerText === csvData[0][1]) { // adjust this based on your CSV data structure
            console.log('Correct!');
        } else {
            console.log('Incorrect!');
        }
        displayNextQuestion();
    }
});