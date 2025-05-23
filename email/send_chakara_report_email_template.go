package sendemail

const ChakaraReportTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Subject}}</title>
	
	<style type="text/css">
	@import url('https://fonts.googleapis.com/css2?family=Afacad:wght@400;500;600;700&display=swap'); /* Added Afacad */
	@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');
	
	</style>
</head>
<body style="margin: 0px;">

	<center class="wrapper" style="width: 100%; table-layout: fixed; background-color: #DAD4CB; padding-bottom: 60px;">

		<table class="main" style="border-spacing: 0; width: 100%; max-width: 666px; background-color: #FDFBF6; font-family: 'Inter', Arial, Helvetica, sans-serif; color: #72513C; text-align: center;" width="100%">

		<tr class="logo-wrapper" style="background-color: #72513C; height: 66px;">
			<td style=""">
				<a href="https://www.ommindshop.com" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/logo.png" alt="Ommind" width="147" height="34" style=""></a>
			</td>
		</tr>
					
		<tr > 
			<td>
				<table align="center" role="presentation" width="100%" style="max-width: 500px;">
					<tr>
						<td align="center">
							<p style="font-weight: 600px; font-size: 30px;">Hi Beautiful Soul,</p>
							<p style="font-size: 14px;">Thank you for taking the OmMind Chakra Test! ✨</p>
							<p style="font-size: 14px;">We’re honored to be part of your journey toward deeper balance, healing, and self-awareness.🧘‍♀️</p>
						</td>
					</tr>
					
				</table>
			</td>
		</tr>

		<tr > 
			<td>
				<table align="center" role="presentation" width="100%" style="max-width: 600px;">
					<tr>
						<td align="center">
							<tr>
								<td height="234" width="601" align="center" valign="middle" background="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/chakara-background.png" style="height: 234px; width: 601px; background-image: url('https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/chakara-background.png'); background-repeat: no-repeat; background-position: center center; background-size: cover; text-align: center; vertical-align: middle;">
									<div>
										<table align="center" role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 500px;"> 
											<tr>
												<td align="center" style="padding: 20px 10px;"> 
													<p style="font-weight: 600; font-size: 16px; color: #ffffff; margin-top: 0; margin-bottom: 10px; font-family: 'Inter', Arial, sans-serif;">Your Chakra Test Number</p>
													<p style="font-size: 14px; color: #ffffff; margin-top: 0; margin-bottom: 10px; font-family: 'Inter', Arial, sans-serif;">{{.ChakaraNumber}}</p>
													<a href="{{.FrontEndUrl}}" style="text-decoration: none; font-size: 14px; color: white; font-family: 'Inter', Arial, sans-serif; padding: 8px 15px; /* Added horizontal padding for better look */ font-weight: 800; background-color: #F8C63E; border-radius: 15px; /* Added border-radius */ margin-top: 10px; /* Added margin-top */ display: inline-block; /* Helps margin-top work reliably */">
														Click to view your chakra results 
														<img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/arrow.png" alt="arrow" width="16" height="15" style="width: 16px; height: 15px; vertical-align: middle; margin-left: 5px; display: inline-block;"> 
													</a>
													
												</td>
											</tr>
										</table>
									</div>
								</td>
							</tr>
						</td>
					</tr>
					
				</table>
			</td>
		</tr>

		
		<tr >
			<td>
				<table align="center" role="presentation" width="100%" style="max-width: 500px; ">
					<tr>
						<td align="center">
							<p style="font-size: 14px;">Your results reveal a unique chakra alignment profile—your personal energy blueprint. Based on your results, we recommend exploring our <strong>customized handmade chakra bracelets</strong>, designed to support your specific energy centers.</p>
                            <p style="font-size: 14px;">You’ve successfully joined our <strong>Chakra Bracelet Giveaway!</strong>  The winner will be revealed on <strong>May 31 at 8 PM (London time)</strong> in our social media stories — stay tuned!</p>
							<p style="margin-top: 0px; font-size: 18px;">As a gift, here’s your exclusive <strong>50% OFF</strong> code:</p>
							<p  style="text-decoration: none;margin-top: 0px; font-size: 14px; color: white; font-family: 'Inter', Arial, sans-serif; padding: 8px 15px; /* Added horizontal padding for better look */ font-weight: 800; background-color: #72513C; border-radius: 15px; /* Added border-radius */ margin-top: 10px; /* Added margin-top */ display: inline-block; /* Helps margin-top work reliably */">
								{{.DiscountCode}}
							</p>
							<p style="font-size: 14px; margin-bottom: 0px;">(valid on one item until June 8, 2025 at 11:59 PM)</p>
                            <p style="font-size: 14px;">Use it to find a bracelet that resonates with your energy and supports your journey.</p>

						</td>
					</tr>
					
				</table>
			</td>
		</tr>

		<tr style="background-color: #72513C;">
			<td>
			  <table width="100%" style="margin-top: 20px;">
				<tr>
				  <td>
					<p style="font-weight: 600; color: #FFFFFF; font-size: 24px; font-family: Inter, Arial, Helvetica, sans-serif; margin: 0; text-align: center;">Discover Your Customized Chakra Bracelet</p>
				  </td>
				</tr>
				<tr>
				  <td style="font-size: 0px; text-align: center;">
					<table style="display: inline-block; width: 100%; max-width: 125px; vertical-align: top; margin: 5px;">
					  <tr>
						<td>
						  <table>
							<tr>
							  <td>
								<a href="https://www.ommindshop.com/products/blessed-weave-tibetan-knotted-bracelet" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/product-1.png" alt="" width="125"  style="border: 0;"></a>
							  </td>
							</tr>
							<tr>
							  <td style="padding: 3px 0 3px;">
								<a href="https://www.ommindshop.com/products/blessed-weave-tibetan-knotted-bracelet" target="_blank" style="text-decoration: none;"><p style="color: #F8C63E; font-size: 14px; font-family: Afacad, Arial, Helvetica, sans-serif; margin-top: 5px;">Blessed Weave – Tibetan Knotted Bracelet</p></a>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
		  
					<table style="display: inline-block; width: 100%; max-width: 125px; vertical-align: top; margin: 5px;">
					  <tr>
						<td>
						  <table>
							<tr>
							  <td>
								<a href="https://www.ommindshop.com/products/monk-shield-chakra-bracelet-men" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/product-2.png" alt="" width="125"  style="border: 0;"></a>
							  </td>
							</tr>
							<tr>
							  <td style="padding: 3px 0 3px;">
								<a href="https://www.ommindshop.com/products/monk-shield-chakra-bracelet-men" target="_blank style="text-decoration: none;"><p style="color: #F8C63E; font-size: 14px; font-family: Afacad, Arial, Helvetica, sans-serif; margin-top: 5px;">Monk Shield - Customized Chakra Bracelet</p></a>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
		  
					<table style="display: inline-block; width: 100%; max-width: 125px; vertical-align: top; margin: 5px;">
					  <tr>
						<td>
						  <table>
							<tr>
							  <td>
								<a href="https://www.ommindshop.com/products/crystal-force-tailored-chakra-energy-bracelet" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/product-3.png" alt="" width="125"  style="border: 0;"></a>
							  </td>
							</tr>
							<tr>
							  <td style="padding: 3px 0 3px;">
								<a href="https://www.ommindshop.com/products/crystal-force-tailored-chakra-energy-bracelet" target="_blank" style="text-decoration: none;"><p style="color: #F8C63E; font-size: 14px; font-family: Afacad, Arial, Helvetica, sans-serif; margin-top: 5px;">Crystal Force – Tailored Chakra Energy Bracelet</p></a>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
		  
					<table style="display: inline-block; width: 100%; max-width: 125px; vertical-align: top; margin: 5px;">
					  <tr>
						<td>
						  <table>
							<tr>
							  <td>
								<a href="https://www.ommindshop.com/products/energy-pathway-custom-chakra-bracelet" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/product-4.png" alt="" width="125"  style="border: 0;"></a>
							  </td>
							</tr>
							<tr>
							  <td style="padding: 3px 0 3px;">
								<a href="https://www.ommindshop.com/products/energy-pathway-custom-chakra-bracelet" target="_blank" style="text-decoration: none;"><p style="color: #F8C63E; font-size: 14px; font-family: Afacad, Arial, Helvetica, sans-serif; margin-top: 5px;">Energy Pathway – Customized Chakra Healing Bracelet</p></a>
							  </td>
							</tr>
						  </table>
						</td>
					  </tr>
					</table>
				  </td>
				</tr>
			  </table>
			</td>
		  </tr>
		  <tr style="background-color: #72513C;">
			<td style="padding: 20px 0;">
			  <hr style="border: none; border-top: 1px solid #FFFFFF; margin: 0 auto; width: 100%;">
			</td>
		  </tr>
		  

		

		<tr>
			<td style="padding:0">
				<table width="100%" style="background-color: #72513C; border-spacing: 0px;">
					<tr>
						<td>
							<tr class="logo-wrapper">
								<td align="center" style="padding-top: 20px;">
									<a href="https://www.ommindshop.com/" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/logo.png" alt="Ommind" width="147" height="34" style=""></a>
								</td>
							</tr>						
						</td>
					</tr>
					<tr >
						<td align="center" style="padding: 25px 25px; font-size: 0;">
							<a href="https://www.facebook.com/profile.php?id=61573346513645" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-1.png"  alt="Social Media 1" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.tiktok.com/@ommindshop" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-2.png"  alt="Social Media 2" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.instagram.com/ommind_shop/" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-3.png"  alt="Social Media 3" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.youtube.com/@OmMind-Official" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-4.png"  alt="Social Media 4" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
						</td>
					</tr>
					<tr>
						<td align="center" style="padding: 0 25px 50px 25px;"> 
							<p style="font-size: 14px; color: white; margin-bottom: 5px; margin-top: 0; margin-bottom: 10px;">86-90 Paul Street, London, EC2A 4UX, United Kingdom</p>
							<p style="font-size: 14px; color: #BDBDBD; margin-top: 0; margin-bottom: 0;">COPYRIGHT ©2025 OmMind Shop All Rights Reserved</p>
						</td>
					</tr>
					
					
				</table>
			</td>
		</tr>			
		</table> 

	</center> 
</body>
</html>
`
