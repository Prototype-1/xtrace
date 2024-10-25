$(document).ready(function () {
    let originalAmount = 0;
    let discountedAmount = 0;

    function showPaymentDetails() {
        $("#paymentDetails").show();
        const paymentType = $("input[name='paymentType']:checked").val();
        if (paymentType === "nol_card_topup") {
            $("#cardTypeContainer").show();
            $("#wallet-balance").hide();
            $("#nol-card-balance").show();
        } else {
            $("#cardTypeContainer").hide();
            $("#wallet-balance").show();
            $("#nol-card-balance").hide();
        }
        fetchCoupon(paymentType);
        $("#userID").trigger("change");
    }

    $("input[name='paymentType']").change(showPaymentDetails);

    $("#userID").on("change", function () {
        const userID = $(this).val();
        const paymentType = $("input[name='paymentType']:checked").val();

        document.getElementById('wallet-balance').innerText = 'Wallet Balance: ₹0';
        document.getElementById('nol-card-balance').innerText = 'Nol Card Balance: ₹0';

        if (userID) {
            if (paymentType !== 'nol_card_topup') {
                fetchWalletBalance(userID);
            }
            fetchNolCardBalance(userID);
        }
    });

    function fetchWalletBalance(userID) {
        fetch(`http://localhost:8000/user/${userID}/wallet/show`)
            .then(response => {
                if (!response.ok) throw new Error('Wallet not found');
                return response.json();
            })
            .then(data => {
                if (data.balance) {
                    document.getElementById('wallet-balance').innerText = `Wallet Balance: ₹${data.balance}`;
                }
            })
            .catch(error => {
                console.error('Error fetching wallet balance:', error);
                document.getElementById('wallet-balance').innerText = 'Wallet not available';
            });
    }

    function fetchNolCardBalance(userID) {
        fetch(`http://localhost:8000/user/${userID}/nol-card/balance/show`)
            .then(response => {
                if (!response.ok) throw new Error('Nol card not found');
                return response.json();
            })
            .then(data => {
                if (data.balance) {
                    document.getElementById('nol-card-balance').innerText = `Nol Card Balance: ₹${data.balance}`;
                }
            })
            .catch(error => {
                console.error('Error fetching Nol card balance:', error);
                document.getElementById('nol-card-balance').innerText = 'Nol card not available';
            });
    }

    $("#id").on("change", function () {
        const paymentType = $("input[name='paymentType']:checked").val();
        const correspondingID = $(this).val();
        console.log("Corresponding ID:", correspondingID);
        if (correspondingID && paymentType) {
            fetchAmount(paymentType, correspondingID);
        }
    });

    function fetchAmount(paymentType, correspondingID) {
        $.ajax({
            url: `http://localhost:8000/user/payment/${paymentType}/${correspondingID}/amount`,
            method: "GET",
            success: function (response) {
                originalAmount = response.amount;
                $("#amount").val(originalAmount);
            },
            error: function (error) {
                $("#paymentResponse").text("Error fetching amount: " + error.responseJSON.error);
            }
        });
    }

    function fetchCoupon(paymentType) {
        $.ajax({
            url: `http://localhost:8000/coupons/${paymentType}`,
            method: "GET",
            success: function (response) {
                if (response && response.coupons && response.coupons.length > 0) {
                    $("#couponDetails").show();
                    let couponDetailsHTML = response.coupons.map(coupon => `
                        <div>
                            <strong>Coupon Code:</strong> ${coupon.code} <br>
                            <strong>Discount:</strong> ${coupon.discount_amount} (${coupon.discount_type})
                        </div>
                        <hr>
                    `).join('');
                    $("#couponDetails").html(couponDetailsHTML);
                } else {
                    $("#couponDetails").hide();
                }
            },
            error: function (error) {
                console.error("Error fetching coupons:", error);
                $("#couponDetails").hide();
            }
        });
    }

    $("#applyCouponButton").click(function () {
        const couponCode = $("#couponCode").val();
        if (couponCode) {
            applyCoupon(couponCode, originalAmount);
        } else {
            $("#paymentResponse").text(`Original Amount: ₹${originalAmount}`);
        }
    });

    function applyCoupon(couponCode, amount) {
        $.ajax({
            url: `http://localhost:8000/coupons/apply`,
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({ coupon_code: couponCode, amount: amount }),
            success: function (response) {
                const discountAmount = response.discount_amount;
                discountedAmount = amount - discountAmount;
                $("#amount").val(discountedAmount.toFixed(2));
                $("#paymentResponse").html(`Coupon applied! Discount: ₹${discountAmount}. Final Amount: ₹${discountedAmount}`);
                $("#createPaymentForm").data("coupon_code", couponCode);
            },
            error: function (error) {
                $("#paymentResponse").html(`Error applying coupon: ${error.responseJSON.error}`);
            }
        });
    }

    $("#createPaymentForm").submit(function (event) {
        event.preventDefault();
        const userID = $("#userID").val();
        const amount = parseFloat($("#amount").val());
        const currency = $("#currency").val();
        const paymentType = $("input[name='paymentType']:checked").val();
        const correspondingID = parseInt($("#id").val(), 10);
        const couponCode = $(this).data("coupon_code");

        if (!userID || isNaN(amount) || !currency || !paymentType || isNaN(correspondingID)) {
            $("#paymentResponse").text("Please fill in all required fields.").css("color", "red");
            return;
        }

        if (amount <= 0) {
            $("#paymentResponse").text("Amount must be greater than ₹0.").css("color", "red");
            return;
        }

        let requestData = {
            amount,
            currency,
            payment_type: paymentType,
            coupon_code: couponCode,
        };

        if (paymentType === "nol_card_topup") {
            const cardType = $("#cardType").val();
            if (!cardType) {
                $("#paymentResponse").text("Please select a card type.").css("color", "red");
                return;
            }
            const minTopup = cardType.toLowerCase() === "gold" ? 100.0 : cardType.toLowerCase() === "silver" ? 50.0 : 20.0;
            if (amount < minTopup) {
                $("#paymentResponse").text(`Minimum top-up for ${cardType} card is ₹${minTopup}.`).css("color", "red");
                return;
            }
            requestData.card_type = cardType;
        }

        if (paymentType === "wallet_topup") {
            requestData.wallet_id = parseInt($("#id").val(), 10);
        } else if (paymentType === "nol_card_topup") {
            requestData.nol_card_id = parseInt($("#id").val(), 10);
        } else if (paymentType === "subscription") {
            requestData.subscription_id = parseInt($("#id").val(), 10);
        } else if (paymentType === "booking") {
            requestData.booking_id = parseInt($("#id").val(), 10);
        }

        requestData[`${paymentType.split('_')[0]}_id`] = correspondingID;

        createPayment(userID, requestData);
    });

    function createPayment(userID, requestData) {
        $.ajax({
            url: `http://localhost:8000/user/${userID}/payment/create`,
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify(requestData),
            success: function (response) {
                console.log('Order ID from backend after creation:', response.order_id);
                $("#paymentResponse").html(`
                    <div style="color: green; font-weight: bold;">
                        Order Creation Successful!
                        <br>Payment Type: ${requestData.payment_type}
                        <br>Order ID: ${response.order_id}
                        <br>Original Amount: ₹${response.original_amount}
                        <br>Discounted Amount: ₹${response.discounted_amount}
                    </div>
                `);
                document.getElementById('orderIdInput').value = response.order_id;
                console.log("Order ID being passed to Razorpay:", response.order_id);
                initiatePayment(response.order_id, response.discounted_amount);
            },
            error: function (error) {
                const errorMessage = error.responseJSON?.error || "Unknown error";
                $("#paymentResponse").text("Error creating payment: " + errorMessage + ". A refund has been initiated.");
                console.log(error);
            }
        });
    }

    function initiatePayment(orderId, amount) {
        var options = {
            key: 'rzp_test_GKYBvaOYPHVdK1',
            name: 'X\' Trace',
            description: 'Payment for Order',
            order_id: orderId,
            amount: amount * 100,
            handler: function (response) {
                console.log("Order ID being passed to Razorpay:", orderId);
                const requestData = {
                    order_id: orderId,
                    payment_id: response.razorpay_payment_id,
                    razorpay_signature: response.razorpay_signature
                };
                verifyPayment(requestData);
            },
            theme: {
                color: '#0000A5'
            }
        };

        var rzp1 = new Razorpay(options);
        rzp1.on('payment.failed', function (response) {
            console.log('Payment failed with error: ', response.error);
            alert(`Payment failed: ${response.error.description}\nCode: ${response.error.code}\nReason: ${response.error.reason}`);
            console.log('Error Metadata:', response.error.metadata);
        });

        rzp1.open();
    }

function checkUserStatus(userId) {
    return $.ajax({
        url: `http://localhost:8000/${userId}/status`,
        method: "GET",
        contentType: "application/json"
    });
}

function disableAllButtons() {
    $('button, input[type="submit"]').prop('disabled', true).addClass('disabled-button');
}

function enableAllButtons() {
    $('button, input[type="submit"]').prop('disabled', false).removeClass('disabled-button');
}

let userBlocked = false;

$("#userID").on('change', function() {
    const userId = $(this).val();
    if (userId) {
        checkUserStatus(userId)
            .done(function(response) {
                userBlocked = response.blocked || response.inactive;
                if (userBlocked) {
                    alert("Your account is currently blocked or inactive. Please contact support for assistance.");
                    disableAllButtons();
                } else {
                    enableAllButtons();
                }
            })
            .fail(function(error) {
                console.log("Error checking user status:", error);
                alert("Error checking user status. Please try again.");
                userBlocked = false;
                enableAllButtons();
            });
    }
});

function verifyPayment(requestData) {
    $.ajax({
        url: "http://localhost:8000/user/payment/verify",
        method: "POST",
        contentType: "application/json",
        data: JSON.stringify({
            order_id: requestData.order_id,
            payment_id: requestData.payment_id,
            razorpay_signature: requestData.razorpay_signature
        }),
        success: function (response) {
            console.log("Payment verification response:", response);
           
            if (response.verified === true) {
                alert(response.message);
                const userID = $("#userID").val();
                const amount = parseFloat($("#amount").val());
                const correspondingID = parseInt($("#id").val(), 10);
               
                switch(response.payment_type) {
                    case 'nol_card_topup':
                        addTopup(userID, amount, correspondingID);
                        break;
                    case 'wallet_topup':
                        alert("Wallet top-up successful!");
                        fetchWalletBalance(userID);
                        break;
                    case 'subscription':
                        alert("Subscription payment successful!");
                        // Subscription-specific actions here if needed
                        break;
                    case 'booking':
                        alert("Booking payment successful!");
                        // Booking-specific actions here if needed
                        break;
                    default:
                        console.log('Unknown payment type');
                }
               
                fetchUpdatedBalances(userID);
            } else {
                if (response.error === "User is blocked or inactive") {
                    alert("Your account is currently blocked or inactive. Please contact support for assistance.");
                    disableAllButtons();
                } else {
                    alert(response.error || "Payment verification failed. Please contact support.");
                }
            }
        },
        error: function (xhr, status, error) {
            console.log("Error during payment verification:");
            console.log("Status:", status);
            console.log("Error:", error);
            console.log("Response:", xhr.responseText);
            
            if (xhr.status === 403 && xhr.responseJSON && xhr.responseJSON.error === "User is blocked or inactive") {
                alert("Your account is currently blocked or inactive. Please contact support for assistance.");
                disableAllButtons();
            } else {
                alert("Error verifying payment. Check console for details.");
            }
        }
    });
}

    function addTopup(userID, amount, correspondingID) {
        const cardType = $("#cardType").val();
        $.ajax({
            url: `http://localhost:8000/user/${userID}/nol-card/topup`,
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                nol_card_id: correspondingID,
                amount: amount,
                card_type: cardType
            }),
            success: function (response) {
                console.log('Top-up successful:', response);
                alert(`Top-up successful! Amount: ₹${response.amount}`);
                fetchUpdatedBalances(userID);
            },
            error: function (error) {
                console.error('Error adding top-up:', error);
                alert(`Error adding top-up: ${error.responseJSON?.error || "Unknown error"}`);
            }
        });
    }
});



