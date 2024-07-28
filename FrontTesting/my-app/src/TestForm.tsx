import React, { useState } from 'react';

// Construct a form comprising of Categories, Subcategories, and Sentences, each accompanied by a checkbox.
// At the end of the form, a submit button will allow the user to submit their selections.
// Selecting a Category automatically selects all its associated Subcategories and Sentences.
// Selecting a Subcategory automatically selects all its associated Sentences.
// Selecting a Sentence automatically marks it, as well as its parent Subcategory and Category, as selected.


// Initial data structure
const initialData = [
  {
    id: 1,
    categoryName: "Kitchen",
    isChecked: false,
    subcategories: [
      {
        id: 11,
        subcategoryName: "Demolition",
        isChecked: false,
        tasks: [
          { id: 111, taskDescription: "Contractor will demolish and remove the existing countertops.", isChecked: false },
          { id: 112, taskDescription: "Contractor will demolish the existing flooring material.", isChecked: false },
          { id: 113, taskDescription: "Contractor will demolish and remove the existing backsplash.", isChecked: false },
          { id: 114, taskDescription: "Contractor will remove an the existing electrical fixtures and plumbing fixtures.", isChecked: false },
          { id: 115, taskDescription: "Contractor will remove the existing kitchen cabinets.", isChecked: false },
          { id: 116, taskDescription: "Contractor will disconnect and remove all existing appliances.", isChecked: false }
        ]
      },
      {
        id: 12,
        subcategoryName: "PLUMBING",
        isChecked: false,
        tasks: [
          { id: 121, taskDescription: "Contractor will upgrade all existing plumbing in \
            accordance with the existing kitchen layout and design. Plumbing to be copper pipe horizontal ¾”, \
            vertical ½” in the kitchen area. Plumbing vents and drainage to be ABS rough-in and ready for installation \
            in the kitchen area all to be done in accordance with the existing kitchen layout and design. \
            All plumbing connections for hot and cold lines, vent pipes, and drain pipes to be connected to the nearest existing location.", isChecked: false },
            { id: 122, taskDescription: "All plumbing to be done according to city codes and regulations.", isChecked: false }
        ]
      },
      {
        id: 13,
        subcategoryName: "ELECTRICAL",
        isChecked: false,
        tasks: [
          { id: 131, taskDescription: "Contractor will install up to 6 new 6” LED recessed lights all to be done according to the existing kitchen layout and design.", isChecked: false },
          { id: 132, taskDescription: "Contractor will install up to 2 new light switches that will be a white decor finish all to be done according to the existing kitchen layout and design.", isChecked: false },
          { id: 132, taskDescription: "Contractor will install up to 6 new GFCI outlets all to be done according to the existing kitchen layout and design.", isChecked: false }
        ]
      }
    ]
  },
  {
    id: 2,
    categoryName: "Bathroom",
    isChecked: false,
    subcategories: [
        {
            id: 21,
            subcategoryName: "Demolition",
            isChecked: false,
            tasks: [
              { id: 211, taskDescription: "Contractor will demolish and remove the existing tile from the bathroom floor and remove the existing tub and tub walls. ", isChecked: false },
              { id: 212, taskDescription: "Contractor will remove all existing switches, GFCI outlets, exhaust fan, light fixtures. ", isChecked: false },
              { id: 213, taskDescription: "Contractor will remove all existing plumbing fixtures such as the faucet and sink (faucet and sink to be stowed away for later reinstallation).", isChecked: false },
              { id: 214, taskDescription: "Contractor will remove all existing bathroom accessories such as towel rack, toilet paper dispenser, and mirror (mirror to be stowed away for later reinstallation).", isChecked: false },
              { id: 215, taskDescription: "Contractor will demolish the proposed wall to accommodate for the proposed bathroom extension.", isChecked: false },
              { id: 216, taskDescription: "Contractor will cut the proposed section of the foundation to accommodate for minimum 5’ bathroom.", isChecked: false }
            ]
          },
          {
            id: 23,
            subcategoryName: "FRAMING",
            isChecked: false,
            tasks: [
              { id: 231, taskDescription: "Contractor will frame a new bathroom layout that will include a new walk-in shower pan with a shampoo niche, vanity area, and toilet area. Bathroom to be approximately 5’ x 9’", isChecked: false },
              { id: 232, taskDescription: "Contractor will frame and block off the existing niche for the open shelving where the existing internet modem is located. ", isChecked: false },
              { id: 232, taskDescription: "Contractor will frame a new niche at the proposed wall to accommodate for the open shelving and reinstall the existing shelves. .", isChecked: false }
            ]
          }
    ]
  }
];

function TestForm() {
    const [data, setData] = useState(initialData);

    const handleCategoryChange = (categoryId: number) => {
        const newData = data.map(category => {
            if (category.id === categoryId) {
                const isChecked = !category.isChecked;
                return {
                    ...category, isChecked, subcategories: category.subcategories.map(sub => ({
                        ...sub, isChecked, tasks: sub.tasks.map(task => ({
                            ...task, isChecked
                        }))
                    }))
                }
            }
            return category;
        });
        setData(newData);
    };

    const handleSubcategoryChange = (categoryId: number, subcategoryId: number) => {
        const newData = data.map(category => {
            if (category.id === categoryId) {
                return {
                    ...category, subcategories: category.subcategories.map(sub => {
                        if (sub.id === subcategoryId) {
                            const isChecked = !sub.isChecked;
                            return {
                                ...sub, isChecked, tasks: sub.tasks.map(task => ({
                                    ...task, isChecked
                                }))
                            }
                        }
                        return sub;
                    })
                }
            }
            return category;
        })
        setData(newData);
    }
    
    const handleTaskChange = (categoryId: number, subcategoryId: number, taskId: number) => {
        const newData = data.map(category => {
            if (categoryId === category.id) {
                return {
                    ...category, subcategories: category.subcategories.map(sub => {
                        if (sub.id === subcategoryId) {
                            return {
                                ...sub, tasks: sub.tasks.map(task => {
                                    if (taskId === task.id) {
                                        return {
                                            ...task, isChecked: !task.isChecked
                                        }
                                    }
                                    return task
                                })
                            }
                        }
                        return sub
                    })
                }
            }
            return category
        })
        setData(newData);
    }
    const handleSubmit = () => {
        const filteredData = data.filter(category => category.isChecked)
                                 .map(category => ({
                                     categoryName: category.categoryName,
                                     subcategories: category.subcategories.filter(sub => sub.isChecked)
                                         .map(sub => ({
                                             subcategoryName: sub.subcategoryName,
                                             tasks: sub.tasks.filter(task => ({
                                                    taskDescription: task.taskDescription
                                             }))
                                         }))
                                 }));
        console.log(filteredData);
    };

    return (
        <div>
            {data.map(category => (
                <div key={category.id}>
                    <h1>
                        <input 
                            type="checkbox"
                            checked={category.isChecked}
                            onChange={()=> handleCategoryChange(category.id)}
                        />
                        {category.categoryName}
                    </h1>
                    {category.subcategories.map(subcategory => (
                        <div key={subcategory.id}>
                            <h3>
                                <input 
                                    type="checkbox"
                                    checked={subcategory.isChecked}
                                    onChange={()=> handleSubcategoryChange(category.id, subcategory.id)}
                                    />
                                {subcategory.subcategoryName}
                            </h3>
                                {subcategory.tasks.map(task =>(
                                    <div key={task.id}>
                                        <input 
                                            type="checkbox"
                                            checked={task.isChecked}
                                            onChange={()=> handleTaskChange(category.id, subcategory.id, task.id)}    
                                        />
                                        {task.taskDescription}
                                    </div>
                                ))}
                        </div>
                    ))}
                </div>
            ))}
            <div style={{display: 'flex', justifyContent: 'center', marginTop: '100px' }}>
                <button onClick={handleSubmit}>Submit</button>
            </div>
        </div>
    );
}

export default TestForm;